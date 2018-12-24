package stemcelldiff_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStemcellDiff(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stemcell diff Suite")
}
