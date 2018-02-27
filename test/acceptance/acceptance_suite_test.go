package acceptance

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
	"os"
)

func TestOpsmanSnitch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var (
	pathToOmenCLI string
	err           error
)

var _ = BeforeSuite(func() {
	os.Unsetenv("OPSMAN_USER")
	os.Unsetenv("OPSMAN_PASSWORD")
	os.Unsetenv("OPSMAN_HOSTNAME")
	pathToOmenCLI, err = gexec.Build("github.com/pivotal-cloudops/omen")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
