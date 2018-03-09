package fakes

import (
	"github.com/pivotal-cloudops/omen/internal/manifest"
)

type FakeManifestsLoader struct {
	LoadAllResponseFunc func(status manifest.ProductStatus) (manifest.Manifests, error)
	LoadResponseFunc    func(status manifest.ProductStatus, tileGuids []string) (manifest.Manifests, error)
}

func (f FakeManifestsLoader) LoadAll(status manifest.ProductStatus) (manifest.Manifests, error) {
	return f.LoadAllResponseFunc(status)
}

func (f FakeManifestsLoader) Load(status manifest.ProductStatus, tileGuids []string) (manifest.Manifests, error) {
	return f.LoadResponseFunc(status, tileGuids)
}
