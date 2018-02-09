package applychanges_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApplychanges(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Applychanges Suite")
}
