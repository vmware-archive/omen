package fakes

import (
	"github.com/pivotal-cloudops/omen/internal/manifest"
	"github.com/pivotal-cloudops/omen/internal/common"
)

type FakeManifestsLoader struct {
	LoadAllResponseFunc func(status common.ProductStatus) (manifest.Manifests, error)
	LoadResponseFunc    func(status common.ProductStatus, tileGuids []string) (manifest.Manifests, error)
}

func (f FakeManifestsLoader) LoadAll(status common.ProductStatus) (manifest.Manifests, error) {
	return f.LoadAllResponseFunc(status)
}

func (f FakeManifestsLoader) Load(status common.ProductStatus, tileGuids []string) (manifest.Manifests, error) {
	return f.LoadResponseFunc(status, tileGuids)
}
