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

	BeforeEach(func() {
		postedURL = ""
		postedBody = ""
	})

	It("Applies changes with no diff", func() {
		manifests := manifest.Manifests{}

		mloader := fakes.FakeManifestsLoader{
			StagedResponseFunc: func() (manifest.Manifests, error) {
				return manifests, nil
			},
			DeployedResponseFunc: func() (manifest.Manifests, error) {
				return manifests, nil
			},
		}

		applychanges.Execute(mloader, mockClient, true)

		Expect(postedURL).To(Equal("/api/v0/installations"))
		Expect(postedBody).To(Equal(applychanges.APPLY_CHANGES_BODY))
	})

	It("Applies changes with diff", func() {
		stagedManifests := manifest.Manifests{
			Data: []manifest.Manifest{
				manifest.Manifest{
					Name: "staged",
				},
			},
		}
		deployedManifests := manifest.Manifests{
			Data: []manifest.Manifest{
				manifest.Manifest{
					Name: "deployed",
				},
			},
		}

		mloader := fakes.FakeManifestsLoader{
			StagedResponseFunc: func() (manifest.Manifests, error) {
				return stagedManifests, nil
			},
			DeployedResponseFunc: func() (manifest.Manifests, error) {
				return deployedManifests, nil
			},
		}

		applychanges.Execute(mloader, mockClient, true)

		Expect(postedURL).To(Equal("/api/v0/installations"))
		Expect(postedBody).To(Equal(applychanges.APPLY_CHANGES_BODY))
	})
})
