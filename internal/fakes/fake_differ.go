package fakes

import "github.com/pivotal-cloudops/omen/internal/manifest"

type FakeDiffer struct {
	DiffResponseFunc func() string
}

func (fd FakeDiffer) Diff(_ manifest.Manifests, _ manifest.Manifests) string {
	return fd.DiffResponseFunc()
}
