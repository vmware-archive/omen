package applychanges_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApplyChanges(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Apply Changes Suite")
}
