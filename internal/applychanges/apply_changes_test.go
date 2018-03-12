package applychanges_test

import (
	"github.com/pivotal-cloudops/omen/internal/applychanges"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cloudops/omen/internal/fakes"
	"github.com/pivotal-cloudops/omen/internal/manifest"
	"reflect"
	"errors"
	"io/ioutil"
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

		applychanges.Execute(mloader, mockClient, []string{}, true, defReportPrinter)

		Expect(postedURL).To(Equal("/api/v0/installations"))
		Expect(postedBody).To(ContainSubstring(`"deploy_products": "all"`))
	})

	It("Selectively applies changes to specified products", func() {
		manifests := manifest.Manifests{}
		mloader := fakes.FakeManifestsLoader{
			LoadResponseFunc: func(status manifest.ProductStatus, tileGuids []string) (manifest.Manifests, error) {
				return manifests, nil
			},
		}

		applychanges.Execute(mloader, mockClient, []string{"product1", "product2"}, true, defReportPrinter)

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

		applychanges.Execute(mloader, mockClient, []string{}, true, defReportPrinter)

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

		applychanges.Execute(mloader, mockClient, []string{}, true, defReportPrinter)

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

		applychanges.Execute(mloader, mockClient, []string{}, true, rp)

		Expect(diff).To(Equal("-manifests.deployed.name=deployed\n+manifests.staged.name=staged\n"))
	})

	Describe("selective tile deployments", func() {
		It("prints out the diff for only the tiles being deployed", func() {

			mloader := fakes.FakeManifestsLoader{

				LoadAllResponseFunc: func(status manifest.ProductStatus) (manifest.Manifests, error) {
					return manifest.Manifests{}, errors.New("loadAll should not be called")
				},

				LoadResponseFunc: func(status manifest.ProductStatus, tileGuids []string) (manifest.Manifests, error) {
					if status == manifest.DEPLOYED && reflect.DeepEqual(tileGuids, []string{"product1", "product2"}) {
						return manifest.Manifests{
							Data: []manifest.Manifest{
								{
									Name: "product1",
								},
								{
									Name: "product2",
								},
							},
						}, nil
					}

					if status == manifest.STAGED && reflect.DeepEqual(tileGuids, []string{"product1", "product2"}) {
						return manifest.Manifests{
							Data: []manifest.Manifest{
								{
									Name: "staged-product1",
								},
								{
									Name: "staged-product2",
								},
							},
						}, nil
					}

					return manifest.Manifests{}, errors.New("don't know how to load these manifests")
				},
			}

			var diff string
			var err error
			rp := fakes.FakeReportPrinter{
				FakeReportFunc: func(s string, e error) {
					diff = s
					err = e
				},
			}

			applychanges.Execute(mloader, mockClient, []string{"product1", "product2"}, true, rp)

			Expect(err).ToNot(HaveOccurred())

			expectedDiff, err := ioutil.ReadFile("testdata/diff.txt")
			Expect(err).ToNot(HaveOccurred())

			Expect(diff).To(Equal(string(expectedDiff)))

		})
	})

})
