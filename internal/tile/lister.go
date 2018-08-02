package tile

import "fmt"

type Lister struct {
	loader tilesLoader
	ui     tableReporter
}

const reportHeader = "Name\tGUID\tVersion\n----\t----\t-------\n"

//go:generate counterfeiter . tilesLoader
type tilesLoader interface {
	LoadDeployed(bool) (Tiles, error)
}

//go:generate counterfeiter . tableReporter
type tableReporter interface {
	Write([]byte) (int, error)
	Flush() error
}

func NewTileLister(tl tilesLoader, ui tableReporter) Lister {
	return Lister{loader: tl, ui: ui}
}

func (l Lister) Execute() error {
	tiles, err := l.loader.LoadDeployed(false)
	if err != nil {
		return err
	}
	if len(tiles.Data) > 0 {
		l.ui.Write([]byte(reportHeader))
		for _, tile := range tiles.Data {
			l.ui.Write([]byte(fmt.Sprintf("%s\t%s\t%s\n",
				tile.Type, tile.GUID, tile.ProductVersion)))
		}
	} else {
		l.ui.Write([]byte("No tiles are installed\n"))
	}
	l.ui.Flush()
	return nil
}
