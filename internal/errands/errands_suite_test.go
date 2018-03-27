package errands_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestErrands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errands Suite")
}
