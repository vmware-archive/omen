package tileguid_test

import (
	"encoding/json"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/omen/internal/fakes"
	"github.com/pivotal-cloudops/omen/internal/tile"
	"github.com/pivotal-cloudops/omen/internal/tileguid"
	"io/ioutil"
)

var _ = Describe("Guid finder", func()	 {
	Context("non-existing product", func() {
		It("returns expected error", func() {
			guid, err := tileguid.FindGuid(setupFakeTilesLoader(), "whatever")

			Expect(guid).To(BeEmpty())
			Expect(err).To(MatchError("product whatever not found"))
		})
	})

	Context("existing product", func() {
		It("returns the guid of the product and no error", func() {
			guid, err := tileguid.FindGuid(setupFakeTilesLoader(), "elastic-runtime")
			Expect(guid).To(Equal("cf-4f9edbd1992fd81250e5"))
			Expect(err).To(Not(HaveOccurred()))
		})
	})

	Context("when the loader fails to load the tiles", func() {
		It("returns the error from the loader", func() {
			tilesLoader := fakes.FakeTilesLoader{}
			tilesLoader.DeployedResponseFunc = func(b bool) (tile.Tiles, error) {
				return tile.Tiles{}, errors.New("Network error")
			}
			guid, err := tileguid.FindGuid(tilesLoader, "whatever")
			Expect(guid).To(BeEmpty())
			Expect(err).To(MatchError("Network error"))
		})
	})
})

func setupFakeTilesLoader() fakes.FakeTilesLoader {
	tilesLoader := fakes.FakeTilesLoader{}
	tilesLoader.DeployedResponseFunc = func(b bool) (tile.Tiles, error) {
		return loadTiles()
	}
	return tilesLoader
}

func loadTiles() (tile.Tiles, error) {

	b, err := ioutil.ReadFile("testdata/tiles.json")
	if err != nil {
		return tile.Tiles{}, err
	}

	var data []*tile.Tile
	err = json.Unmarshal(b, &data)
	if err != nil {
		return tile.Tiles{}, err
	}

	return tile.Tiles{
		Data: data,
	}, err

}
