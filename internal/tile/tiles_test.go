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
	Context("#Write", func() {
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
})
