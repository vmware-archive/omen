package fakes

import (
	"github.com/pivotal-cloudops/omen/internal/tile"
)

type FakeTilesLoader struct {
	StagedResponseFunc   func(bool) (tile.Tiles, error)
	DeployedResponseFunc func(bool) (tile.Tiles, error)
}

func (f FakeTilesLoader) LoadStaged(fetchTileMetadata bool) (tile.Tiles, error) {
	return f.StagedResponseFunc(fetchTileMetadata)
}

func (f FakeTilesLoader) LoadDeployed(fetchTileMetadata bool) (tile.Tiles, error) {
	return f.DeployedResponseFunc(fetchTileMetadata)
}
