package applychanges_test

import (
	"github.com/pivotal-cloudops/omen/internal/applychanges"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cloudops/omen/internal/fakes"
	"github.com/pivotal-cloudops/omen/internal/manifest"
)

var _ = Describe("Apply Changes - Execute", func() {
	var postedURL, postedBody string
	mockClient := fakes.FakeOMClient{
		PostFunc: func(url string, body string) ([]byte, error) {
			postedURL = url
			postedBody = body

			return nil, nil
		},
	}

	defReportPrinter := fakes.FakeReportPrinter{
		FakeReportFunc: func(string, error) {
			return
		},
	}

	BeforeEach(func() {
		postedURL = ""
		postedBody = ""
	})

	It("Applies all changes by default", func() {
		manifests := manifest.Manifests{}

		mloader := fakes.FakeManifestsLoader{
			LoadAllResponseFunc: func(status manifest.ProductStatus) (manifest.Manifests, error) {
				return manifests, nil
			},
		}

		applychanges.Execute(mloader, mockClient, "", true, defReportPrinter)

		Expect(postedURL).To(Equal("/api/v0/installations"))
		Expect(postedBody).To(ContainSubstring(`"deploy_products": "all"`))
	})

	It("Selectively applies changes to specified products", func() {
		manifests := manifest.Manifests{}
		mloader := fakes.FakeManifestsLoader{
			LoadAllResponseFunc: func(status manifest.ProductStatus) (manifest.Manifests, error) {
				return manifests, nil
			},
		}

		applychanges.Execute(mloader, mockClient, "product1,product2", true, defReportPrinter)

		Expect(postedURL).To(Equal("/api/v0/installations"))
		Expect(postedBody).To(ContainSubstring(`"deploy_products": "product1,product2"`))
	})

	It("Applies changes with no difference between staged and deployed", func() {
		manifests := manifest.Manifests{}

		mloader := fakes.FakeManifestsLoader{
			LoadAllResponseFunc: func(status manifest.ProductStatus) (manifest.Manifests, error) {
				return manifests, nil
			},
		}

		applychanges.Execute(mloader, mockClient, "", true, defReportPrinter)

		Expect(postedURL).To(Equal("/api/v0/installations"))
		Expect(postedBody).To(ContainSubstring(`"deploy_products": "all"`))
	})

	It("Applies changes with difference between staged and deployed", func() {
		stagedManifests := manifest.Manifests{
			Data: []manifest.Manifest{
				{
					Name: "staged",
				},
			},
		}
		deployedManifests := manifest.Manifests{
			Data: []manifest.Manifest{
				{
					Name: "deployed",
				},
			},
		}

		mloader := fakes.FakeManifestsLoader{
			LoadAllResponseFunc: func(status manifest.ProductStatus) (manifest.Manifests, error) {
				if status == manifest.DEPLOYED {
					return deployedManifests, nil
				}
				return stagedManifests, nil
			},
		}

		applychanges.Execute(mloader, mockClient, "", true, defReportPrinter)

		Expect(postedURL).To(Equal("/api/v0/installations"))
		Expect(postedBody).To(ContainSubstring(`"deploy_products": "all"`))
	})

	It("Prints out the the diff between all staged and deployed tiles", func() {
		stagedManifests := manifest.Manifests{
			Data: []manifest.Manifest{
				{
					Name: "staged",
				},
			},
		}
		deployedManifests := manifest.Manifests{
			Data: []manifest.Manifest{
				{
					Name: "deployed",
				},
			},
		}

		mloader := fakes.FakeManifestsLoader{
			LoadAllResponseFunc: func(status manifest.ProductStatus) (manifest.Manifests, error) {
				if status == manifest.DEPLOYED {
					return deployedManifests, nil
				}
				return stagedManifests, nil
			},
		}

		var diff string
		rp := fakes.FakeReportPrinter{
			FakeReportFunc: func(s string, e error) {
				diff = s
			},
		}

		applychanges.Execute(mloader, mockClient, "", true, rp)

		Expect(diff).To(Equal("-manifests.deployed.name=deployed\n+manifests.staged.name=staged\n"))
	})
})
