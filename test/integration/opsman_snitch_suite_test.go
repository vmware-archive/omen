package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestOpsmanSnitch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var (
	pathToOmenCLI string
	err           error
)

var _ = BeforeSuite(func() {
	pathToOmenCLI, err = gexec.Build("github.com/pivotal-cloudops/omen")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
