package tileguid

import (
	"github.com/pivotal-cloudops/omen/internal/tile"
)

type tileLoader interface {
	LoadDeployed(fetchTileMetadata bool) (tile.Tiles, error)
}

func FindGuid(tileLoader tileLoader, productSlug string) (string, error) {
	tiles, err := tileLoader.LoadDeployed(false)
	if err != nil {
		return "", err
	}

	foundTile, err := tiles.FindBySlug(productSlug)
	if err != nil {
		return "", err
	}

	return foundTile.GUID, nil
}
