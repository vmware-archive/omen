package tile_test

import (
	"encoding/json"
	"io/ioutil"

	"os"

	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/omen/internal/tile"
)

var _ = Describe("Tiles", func() {
	Describe("#Write", func() {
		It("writes out all the Tiles at the specified path", func() {
			var t tile.Tile
			f, err := ioutil.ReadFile("testdata/expected_tile.json")
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal(f, &t)
			Expect(err).NotTo(HaveOccurred())

			tiles := tile.Tiles{Data: []*tile.Tile{&t}}

			tmpdir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpdir)

			tiles.Write(tmpdir)

			dir, err := os.Open(tmpdir + "/cf")
			Expect(err).NotTo(HaveOccurred())
			defer dir.Close()

			filenames, err := dir.Readdirnames(-1)
			Expect(err).NotTo(HaveOccurred())

			sort.Slice(filenames, func(i, j int) bool {
				a := filenames[i]
				b := filenames[j]
				return a < b
			})

			Expect(filenames).To(Equal([]string{
				"errands.json",
				"networks_and_azs.json",
				"properties.json",
				"resources.json",
			}))
		})
	})

	Describe("FindBySlug", func() {
		It("it finds the guid for the specified slug", func() {
			t := tile.Tiles {
				Data: []*tile.Tile {
					{
						GUID: "tile-1234",
						Type: "tile",
					},
				},
			}

			ts, err := t.FindBySlug("tile")

			Expect(err).To(Not(HaveOccurred()))
			Expect(ts.GUID).To(Equal("tile-1234"))
		})

		It("returns an error when the tile is not found", func() {
			t := tile.Tiles {
				Data: []*tile.Tile {
					{
						GUID: "tile-1234",
						Type: "tile",
					},
				},
			}

			_, err := t.FindBySlug("tile-that-does-not-exist")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("product tile-that-does-not-exist not found"))
		})
	})
})
