package tile_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/kylelemons/godebug/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/omen/internal/fakes"
	"github.com/pivotal-cloudops/omen/internal/tile"
)

var _ = Describe("Tiles Loader", func() {

	It("should load deployed tiles", func() {
		fakeOMClient := fakes.FakeOMClient{
			GetFunc: func(endpoint string) ([]byte, error) {
				switch endpoint {
				case "/api/v0/deployed/products":
					return ioutil.ReadFile("testdata/tiles.json")
				default:
					return nil, errors.New(fmt.Sprintf("invalid endpoint %v", endpoint))
				}
			},
		}
		loader := tile.NewTilesLoader(fakeOMClient)
		tiles, err := loader.LoadDeployed(false)
		Expect(err).NotTo(HaveOccurred())
		data := tiles.Data
		Expect(data).To(HaveLen(1))
	})

	Context("loading without metadata", func() {
		It("should load", func() {
			fakeOMClient := fakes.FakeOMClient{
				GetFunc: func(endpoint string) ([]byte, error) {
					switch endpoint {
					case "/api/v0/staged/products":
						return ioutil.ReadFile("testdata/tiles.json")
					default:
						return nil, errors.New(fmt.Sprintf("invalid endpoint %v", endpoint))
					}
				},
			}
			loader := tile.NewTilesLoader(fakeOMClient)
			tiles, err := loader.LoadStaged(false)
			Expect(err).NotTo(HaveOccurred())
			data := tiles.Data
			Expect(data).To(HaveLen(1))
			tile := data[0]

			Expect(tile.InstallationName).To(Equal("cf-4f9edbd1992fd81250e5"))
			Expect(tile.GUID).To(Equal("guid"))
			Expect(tile.Type).To(Equal("cf"))
			Expect(tile.ProductVersion).To(Equal("1.12.0.0"))
			Expect(tile.Networks["networks_and_azs"]).To(BeNil())
			Expect(tile.Errands["errands"]).To(BeNil())
			Expect(tile.Resources["resources"]).To(BeNil())
			Expect(tile.Properties["properties"]).To(BeNil())
		})
	})

	Context("loading with metadata", func() {

		DescribeTable("should load", func(status string) {
			fakeOMClient := fakes.FakeOMClient{
				GetFunc: func(endpoint string) ([]byte, error) {
					switch endpoint {
					case fmt.Sprintf("/api/v0/%s/products", status):
						return ioutil.ReadFile("testdata/tiles.json")
					case fmt.Sprintf("/api/v0/%s/products/guid/networks_and_azs", status):
						return ioutil.ReadFile("testdata/cf/networks_and_azs.json")
					case fmt.Sprintf("/api/v0/%s/products/guid/errands", status):
						return ioutil.ReadFile("testdata/cf/errands.json")
					case fmt.Sprintf("/api/v0/%s/products/guid/resources", status):
						return ioutil.ReadFile("testdata/cf/resources.json")
					case fmt.Sprintf("/api/v0/%s/products/guid/properties", status):
						return ioutil.ReadFile("testdata/cf/properties.json")
					default:
						return nil, errors.New(fmt.Sprintf("invalid endpoint %v", endpoint))
					}
				},
			}

			loader := tile.NewTilesLoader(fakeOMClient)

			var (
				tiles tile.Tiles
				err   error
			)

			switch status {
			case "staged":
				tiles, err = loader.LoadStaged(true)
			case "deployed":
				tiles, err = loader.LoadDeployed(true)
			default:
				err = errors.New("invalid product status")
			}

			Expect(err).NotTo(HaveOccurred())
			data := tiles.Data
			Expect(data).To(HaveLen(1))
			actualTile := data[0]
			var expectedTile tile.Tile
			tileJSON, err := ioutil.ReadFile("testdata/expected_tile.json")
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal(tileJSON, &expectedTile)
			Expect(err).NotTo(HaveOccurred())

			diff := pretty.Compare(actualTile, expectedTile)
			Expect(diff).To(BeEmpty())
		},
			Entry("staged", "staged"),
			Entry("deployed", "deployed"))

	})

	It("should fail if fetching tiles fails", func() {
		fakeOMClient := fakes.FakeOMClient{
			GetFunc: func(endpoint string) ([]byte, error) {
				return nil, errors.New("failed")
			},
		}
		finder := tile.NewTilesLoader(fakeOMClient)
		_, err := finder.LoadStaged(true)
		Expect(err).To(HaveOccurred())
	})

	It("should fail if loading tile metadata fails", func() {
		fakeOMClient := fakes.FakeOMClient{
			GetFunc: func(endpoint string) ([]byte, error) {
				switch endpoint {
				case "/api/v0/staged/products":
					return ioutil.ReadFile("testdata/tiles.json")
				case "/api/v0/staged/products/guid/networks_and_azs":
					return nil, errors.New("loading networks failed")
				case "/api/v0/staged/products/guid/errands":
					return ioutil.ReadFile("testdata/cf/errands.json")
				case "/api/v0/staged/products/guid/resources":
					return ioutil.ReadFile("testdata/cf/resources.json")
				case "/api/v0/staged/products/guid/properties":
					return ioutil.ReadFile("testdata/cf/properties.json")
				default:
					return nil, errors.New(fmt.Sprintf("invalid endpoint %v", endpoint))
				}
			},
		}

		finder := tile.NewTilesLoader(fakeOMClient)
		_, err := finder.LoadStaged(true)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("loading networks failed"))
	})

})
