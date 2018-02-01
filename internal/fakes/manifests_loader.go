package fakes

import (
	"github.com/pivotal-cloudops/omen/internal/manifest"
)

type FakeManifestsLoader struct {
	StagedResponseFunc   func() (manifest.Manifests, error)
	DeployedResponseFunc func() (manifest.Manifests, error)
}

func (f FakeManifestsLoader) LoadStaged() (manifest.Manifests, error) {
	return f.StagedResponseFunc()
}

func (f FakeManifestsLoader) LoadDeployed() (manifest.Manifests, error) {
	return f.DeployedResponseFunc()
}
