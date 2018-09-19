package tileguid_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGuid(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Guid Suite")
}
