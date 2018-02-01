package fakes

import (
	"github.com/pivotal-cloudops/omen/internal/tile"
)

type FakeTilesLoader struct {
	StagedResponseFunc   func() (tile.Tiles, error)
	DeployedResponseFunc func() (tile.Tiles, error)
}

func (f FakeTilesLoader) LoadStaged(_ bool) (tile.Tiles, error) {
	return f.StagedResponseFunc()
}

func (f FakeTilesLoader) LoadDeployed(_ bool) (tile.Tiles, error) {
	return f.DeployedResponseFunc()
}
