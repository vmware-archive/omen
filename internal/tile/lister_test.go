package tile_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/omen/internal/tile"
	"github.com/pivotal-cloudops/omen/internal/tile/tilefakes"
)

var (
	subject  tile.Lister
	loader   *tilefakes.FakeTilesLoader
	reporter *tilefakes.FakeTableReporter
)

var _ = Describe("Lister", func() {

	mockTile := &tile.Tile{
		InstallationName: "⸘why am I totally random‽",
		Type:             "spinner",
		GUID:             "spinner-2b35f5d3fd3ed898a798d79b",
		ProductVersion:   "12.34.56",
	}

	BeforeEach(func() {
		loader = &tilefakes.FakeTilesLoader{}
		reporter = &tilefakes.FakeTableReporter{}
		subject = tile.NewTileLister(loader, reporter)
	})

	It("is deleted", func() {
		Expect(true).To(BeTrue())
	})

	It("prints an empty report", func() {
		err := subject.Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(reporter.WriteCallCount()).To(Equal(1))
		line := string(reporter.WriteArgsForCall(0))
		Expect(line).To(HavePrefix("No tiles are installed\n"))
		Expect(reporter.FlushCallCount()).To(Equal(1))
	})

	It("prints a tile report", func() {
		tiles := tile.Tiles{
			Data: []*tile.Tile{mockTile},
		}
		loader.LoadDeployedReturns(tiles, nil)

		err := subject.Execute()
		Expect(err).NotTo(HaveOccurred())

		Expect(reporter.WriteCallCount()).To(Equal(2))

		header := string(reporter.WriteArgsForCall(0))
		Expect(header).To(Equal("Name\tGUID\tVersion\n----\t----\t-------\n"))

		tileLine := string(reporter.WriteArgsForCall(1))
		Expect(tileLine).To(Equal("spinner\tspinner-2b35f5d3fd3ed898a798d79b\t12.34.56\n"))
	})

	It("flushes output at the end", func() {
	    tiles := tile.Tiles{
	    	Data: []*tile.Tile{mockTile, mockTile, mockTile},
		}
	    loader.LoadDeployedReturns(tiles, nil)

	    reporter.FlushStub = func() error {
			Expect(reporter.WriteCallCount()).To(Equal(4))
			return nil
		}

	    err := subject.Execute()
	    Expect(err).NotTo(HaveOccurred())
	    Expect(reporter.FlushCallCount()).To(Equal(1))
	})

	It("surfaces the underlying errors", func() {
		loader.LoadDeployedReturns(tile.Tiles{}, errors.New("boom"))
		err := subject.Execute()

		Expect(err).To(MatchError(errors.New("boom")))
	})
})
