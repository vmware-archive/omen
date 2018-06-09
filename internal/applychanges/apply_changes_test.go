package applychanges_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"io/ioutil"
	"reflect"

	"github.com/pivotal-cloudops/omen/internal/applychanges"
	"github.com/pivotal-cloudops/omen/internal/applychanges/applychangesfakes"
	"github.com/pivotal-cloudops/omen/internal/common"
	"github.com/pivotal-cloudops/omen/internal/fakes"
	"github.com/pivotal-cloudops/omen/internal/manifest"
	"github.com/pivotal-cloudops/omen/internal/tile"
)

var (
	twoTiles = tile.Tiles{
		Data: []*tile.Tile{
			{
				GUID: "guid1",
				Type: "product1",
			},
			{
				GUID: "guid2",
				Type: "product2",
			},
		},
	}
)

var _ = Describe("Apply Changes - Execute", func() {
	var mockClient *applychangesfakes.FakeOpsmanClient
	var reportPrinter *applychangesfakes.FakeReportPrinter

	BeforeEach(func() {
		mockClient = &applychangesfakes.FakeOpsmanClient{}
		reportPrinter = &applychangesfakes.FakeReportPrinter{}
	})

	It("Applies all changes by default", func() {
		manifests := manifest.Manifests{}

		manifestsLoader := &applychangesfakes.FakeManifestsLoader{
			LoadAllStub: func(status common.ProductStatus) (manifest.Manifests, error) {
				return manifests, nil
			},
		}

		tilesLoader := fakes.FakeTilesLoader{}

		subject := applychanges.NewApplyChangesOp(manifestsLoader, tilesLoader, mockClient, reportPrinter, applychanges.ApplyChangesOptions{TileSlugs: []string{}, NonInteractive: true})
		subject.Execute()

		postedUrl, postedBody, _ := mockClient.PostArgsForCall(0)
		Expect(postedUrl).To(Equal("/api/v0/installations"))
		Expect(postedBody).To(ContainSubstring(`"deploy_products": "all"`))
	})

	Describe("no changes between staged and deployed", func() {
		var manifests manifest.Manifests
		var manifestLoader *applychangesfakes.FakeManifestsLoader
		var tilesLoader fakes.FakeTilesLoader
		var subject applychanges.ApplyChangesOp

		BeforeEach(func() {
			manifests = manifest.Manifests{}
			manifestLoader = &applychangesfakes.FakeManifestsLoader{
				LoadAllStub: func(status common.ProductStatus) (manifest.Manifests, error) {
					return manifests, nil
				},
			}
			tilesLoader = fakes.FakeTilesLoader{}
			subject = applychanges.NewApplyChangesOp(
				manifestLoader,
				tilesLoader,
				mockClient,
				reportPrinter,
				applychanges.ApplyChangesOptions{TileSlugs: []string{}, NonInteractive: true})
		})

		It("applies changes", func() {
			subject.Execute()

			postedUrl, postedBody, _ := mockClient.PostArgsForCall(0)
			Expect(postedUrl).To(Equal("/api/v0/installations"))
			Expect(postedBody).To(ContainSubstring(`"deploy_products": "all"`))
		})

		It("produces a warning for a full run", func() {
			subject.Execute()

			Expect(reportPrinter.Invocations()).To(HaveLen(1))
			warning := reportPrinter.PrintReportArgsForCall(0)
			Expect(warning).To(ContainSubstring("Warning:"))
			Expect(warning).To(ContainSubstring("no pending changes"))
		})
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

		tilesLoader := fakes.FakeTilesLoader{}

		manifestsLoader := &applychangesfakes.FakeManifestsLoader{
			LoadAllStub: func(status common.ProductStatus) (manifest.Manifests, error) {
				if status == common.DEPLOYED {
					return deployedManifests, nil
				}
				return stagedManifests, nil
			},
		}

		subject := applychanges.NewApplyChangesOp(manifestsLoader, tilesLoader, mockClient, reportPrinter, applychanges.ApplyChangesOptions{TileSlugs: []string{}, NonInteractive: true})
		subject.Execute()

		postedUrl, postedBody, _ := mockClient.PostArgsForCall(0)
		Expect(postedUrl).To(Equal("/api/v0/installations"))
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

		manifestsLoader := &applychangesfakes.FakeManifestsLoader{
			LoadAllStub: func(status common.ProductStatus) (manifest.Manifests, error) {
				if status == common.DEPLOYED {
					return deployedManifests, nil
				}
				return stagedManifests, nil
			},
		}

		tilesLoader := fakes.FakeTilesLoader{}

		subject := applychanges.NewApplyChangesOp(manifestsLoader, tilesLoader, mockClient, reportPrinter, applychanges.ApplyChangesOptions{TileSlugs: []string{}, NonInteractive: true})
		subject.Execute()
		diff := reportPrinter.PrintReportArgsForCall(0)
		Expect(diff).To(Equal("-manifests.deployed.name=deployed\n+manifests.staged.name=staged\n"))
	})

	Describe("selective tile deployments", func() {
		It("applies changes to specified products", func() {
			fetchTileMetadata := true
			manifests := manifest.Manifests{}
			manifestsLoader := &applychangesfakes.FakeManifestsLoader{
				LoadStub: func(status common.ProductStatus, tileGuids []string) (manifest.Manifests, error) {
					return manifests, nil
				},
			}

			tilesLoader := fakes.FakeTilesLoader{
				StagedResponseFunc: func(b bool) (tile.Tiles, error) {
					fetchTileMetadata = b
					return twoTiles, nil
				},
			}

			subject := applychanges.NewApplyChangesOp(manifestsLoader, tilesLoader, mockClient, reportPrinter, applychanges.ApplyChangesOptions{TileSlugs: []string{"product1", "product2"}, NonInteractive: true})
			subject.Execute()

			Expect(fetchTileMetadata).To(BeFalse())

			postedUrl, postedBody, _ := mockClient.PostArgsForCall(0)
			Expect(postedUrl).To(Equal("/api/v0/installations"))
			Expect(postedBody).To(ContainSubstring(`"deploy_products": "guid1,guid2"`))
		})

		It("fails when slug not found", func() {
			manifestsLoader := &applychangesfakes.FakeManifestsLoader{
				LoadStub: func(status common.ProductStatus, tileGuids []string) (manifest.Manifests, error) {
					return manifest.Manifests{}, nil
				},
			}

			tilesLoader := fakes.FakeTilesLoader{
				StagedResponseFunc: func(b bool) (tile.Tiles, error) {
					return twoTiles, nil
				},
			}

			subject := applychanges.NewApplyChangesOp(manifestsLoader, tilesLoader, mockClient, reportPrinter, applychanges.ApplyChangesOptions{TileSlugs: []string{"product3", "product2"}, NonInteractive: true})
			err := subject.Execute()

			Expect(err).To(HaveOccurred())

			Expect(mockClient.PostCallCount()).To(BeZero())
		})

		It("fails when tile loading fails", func() {
			manifestsLoader := &applychangesfakes.FakeManifestsLoader{
				LoadStub: func(status common.ProductStatus, tileGuids []string) (manifest.Manifests, error) {
					return manifest.Manifests{}, nil
				},
			}

			tilesLoader := fakes.FakeTilesLoader{
				StagedResponseFunc: func(b bool) (tile.Tiles, error) {
					return tile.Tiles{}, errors.New("can't load tiles")
				},
			}

			subject := applychanges.NewApplyChangesOp(manifestsLoader, tilesLoader, mockClient, reportPrinter, applychanges.ApplyChangesOptions{TileSlugs: []string{"product3"}, NonInteractive: true})
			err := subject.Execute()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("can't load tiles"))
			Expect(mockClient.PostCallCount()).To(BeZero())
		})

		It("prints out the diff for only the tiles being deployed", func() {

			manifestsLoader := &applychangesfakes.FakeManifestsLoader{

				LoadAllStub: func(status common.ProductStatus) (manifest.Manifests, error) {
					return manifest.Manifests{}, errors.New("loadAll should not be called")
				},

				LoadStub: func(status common.ProductStatus, tileGuids []string) (manifest.Manifests, error) {
					if status == common.DEPLOYED && reflect.DeepEqual(tileGuids, []string{"guid1", "guid2"}) {
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

					if status == common.STAGED && reflect.DeepEqual(tileGuids, []string{"guid1", "guid2"}) {
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

			tilesLoader := fakes.FakeTilesLoader{
				StagedResponseFunc: func(b bool) (tile.Tiles, error) {
					return twoTiles, nil
				},
			}

			subject := applychanges.NewApplyChangesOp(manifestsLoader, tilesLoader, mockClient, reportPrinter, applychanges.ApplyChangesOptions{TileSlugs: []string{"product1", "product2"}, NonInteractive: true})
			subject.Execute()
			diff := reportPrinter.PrintReportArgsForCall(0)

			expectedDiff, err := ioutil.ReadFile("testdata/diff.txt")
			Expect(err).ToNot(HaveOccurred())

			Expect(diff).To(Equal(string(expectedDiff)))

		})
	})

	Describe("Dry run", func() {
		It("it outputs the diff but does not apply changes", func() {
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

			manifestsLoader := &applychangesfakes.FakeManifestsLoader{
				LoadAllStub: func(status common.ProductStatus) (manifest.Manifests, error) {
					if status == common.DEPLOYED {
						return deployedManifests, nil
					}
					return stagedManifests, nil
				},
			}

			tilesLoader := fakes.FakeTilesLoader{}

			subject := applychanges.NewApplyChangesOp(manifestsLoader, tilesLoader, mockClient, reportPrinter, applychanges.ApplyChangesOptions{[]string{}, true, true, false})
			subject.Execute()
			diff := reportPrinter.PrintReportArgsForCall(0)
			Expect(diff).To(Equal("-manifests.deployed.name=deployed\n+manifests.staged.name=staged\n"))

			Expect(mockClient.PostCallCount()).To(BeZero())
		})
	})

	Describe("Quiet run", func() {
		It("only outputs the return of ops manager apply changes", func() {
			applyChangesReply := `{"install":{"id": 303}}`

			tilesLoader := fakes.FakeTilesLoader{}
			manifestsLoader := &applychangesfakes.FakeManifestsLoader{}

			mockClient.PostReturns([]byte(applyChangesReply), nil)

			subject := applychanges.NewApplyChangesOp(manifestsLoader, tilesLoader, mockClient, reportPrinter, applychanges.ApplyChangesOptions{[]string{}, true, false, true})
			subject.Execute()
			Expect(reportPrinter.PrintReportCallCount()).To(Equal(1))

			Expect(reportPrinter.PrintReportArgsForCall(0)).To(MatchJSON(applyChangesReply))
		})
	})

})
