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

	Describe("#FindBySlugsOrGUIDs", func() {
		It("returns an empty slice and no error for empty input", func() {
			t := tile.Tiles{}

			ts, err := t.FindBySlugsOrGUIDs([]string{})
			Expect(err).To(Not(HaveOccurred()))
			Expect(ts).To(Equal([]*tile.Tile{}))
		})

		It("returns one element slice for a single slug", func() {
			t := tile.Tiles{
				Data: []*tile.Tile{
					{
						GUID: "tile-1234",
						Type: "tile",
					},
				},
			}

			ts, err := t.FindBySlugsOrGUIDs([]string{"tile"})
			Expect(err).To(Not(HaveOccurred()))
			Expect(ts).To(Equal(t.Data))

		})

		It("returns one element slice for a single GUID", func() {
			t := tile.Tiles{
				Data: []*tile.Tile{
					{
						GUID: "tile-1234",
						Type: "tile",
					},
				},
			}

			ts, err := t.FindBySlugsOrGUIDs([]string{"tile-1234"})
			Expect(err).To(Not(HaveOccurred()))
			Expect(ts).To(Equal(t.Data))
		})

		It("returns a list of tiles when list of slugs is used", func() {
			t := tile.Tiles{
				Data: []*tile.Tile{
					{
						GUID: "tile-1234",
						Type: "tile",
					},
					{
						GUID: "tile2-5678",
						Type: "tile2",
					},
				},
			}

			ts, err := t.FindBySlugsOrGUIDs([]string{"tile", "tile2"})
			Expect(err).To(Not(HaveOccurred()))
			Expect(ts).To(Equal(t.Data))
		})

		It("returns a list of tiles when list of guids is used", func() {
			t := tile.Tiles{
				Data: []*tile.Tile{
					{
						GUID: "tile-1234",
						Type: "tile",
					},
					{
						GUID: "tile2-5678",
						Type: "tile2",
					},
				},
			}

			ts, err := t.FindBySlugsOrGUIDs([]string{"tile-1234", "tile2-5678"})
			Expect(err).To(Not(HaveOccurred()))
			Expect(ts).To(Equal(t.Data))
		})

		It("returns an error when a tile is not found for a given value", func() {
			t := tile.Tiles{
				Data: []*tile.Tile{
					{
						GUID: "tile-1234",
						Type: "tile",
					},
					{
						GUID: "tile2-5678",
						Type: "tile2",
					},
				},
			}

			_, err := t.FindBySlugsOrGUIDs([]string{"tile10-1234", "tile2-5678"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("product tile10-1234 is not found"))
		})

		It("returns an error when input is a mix of slugs and guids", func() {
			t := tile.Tiles{
				Data: []*tile.Tile{
					{
						GUID: "tile-1234",
						Type: "tile",
					},
					{
						GUID: "tile2-5678",
						Type: "tile2",
					},
				},
			}

			_, err := t.FindBySlugsOrGUIDs([]string{"tile", "tile2-5678"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("input contains a mix of GUIDs and names"))
		})
	})

	Describe("#FindBySlug", func() {
		It("finds the guid for the specified slug", func() {
			t := tile.Tiles{
				Data: []*tile.Tile{
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
			t := tile.Tiles{
				Data: []*tile.Tile{
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

	Describe("#FindByGuid", func() {
		It("finds a tile by guid", func() {
			t := tile.Tiles{
				Data: []*tile.Tile{
					{
						GUID: "tile-1234",
						Type: "tile",
					},
				},
			}

			ts, err := t.FindByGuid("tile-1234")

			Expect(err).To(Not(HaveOccurred()))
			Expect(ts.GUID).To(Equal("tile-1234"))
		})

		It("returns an error when a tile is not found", func() {
			t := tile.Tiles{
				Data: []*tile.Tile{
					{
						GUID: "tile-1234",
						Type: "tile",
					},
				},
			}

			_, err := t.FindByGuid("rummage-island")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("product guid rummage-island not found"))
		})
	})
})
