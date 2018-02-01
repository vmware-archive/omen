package acceptance_tests

import (
	"os/exec"

	"strings"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("omen", func() {
	var session *gexec.Session
	validUsername := os.Getenv("OPSMAN_USER")
	opsManagerTimeOut := "120s"

	It("should print diagnostic information", func() {
		command := exec.Command(pathToOmenCLI, "diagnostics")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func() string { return string(session.Out.Contents()) }, opsManagerTimeOut).Should(ContainSubstring("infrastructure_type\": \"google"))
		Eventually(session, opsManagerTimeOut).Should(gexec.Exit(0))
	})

	It("should output staged tiles information", func() {
		defer os.RemoveAll("/tmp/snitch")
		command := exec.Command(pathToOmenCLI, "staged-tiles", "-o=/tmp/snitch")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, opsManagerTimeOut).Should(gexec.Exit(0))
	})

	It("should print manifests report", func() {
		command := exec.Command(pathToOmenCLI, "manifests")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func() string { return string(session.Out.Contents()) }, opsManagerTimeOut).Should(ContainSubstring("p-bosh"))
		Eventually(func() int { return strings.Count(string(session.Out.Contents()), "\"name\": \"p-bosh\"") }, opsManagerTimeOut).Should(BeNumerically(">=", 1))
		Eventually(func() int { return strings.Count(string(session.Out.Contents()), "\"cloud_config\"") }, opsManagerTimeOut).Should(BeNumerically("==", 1))
		Eventually(session, opsManagerTimeOut).Should(gexec.Exit(0))
	})

	It("should not error on apply changes", func() {
		command := exec.Command(pathToOmenCLI, "apply-changes", "-n=true")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, opsManagerTimeOut).Should(gexec.Exit(0))
	})

	It("should print an error and exit nonzero with bad credentials", func() {
		command := exec.Command(pathToOmenCLI, "-u=foo", "diagnostics")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func() string { return string(session.Out.Contents()) }, opsManagerTimeOut).Should(ContainSubstring("Bad credentials"))
		Eventually(session, opsManagerTimeOut).Should(Not(gexec.Exit(0)))
	})

	It("should print an error and exit with nonzero value if no command parameter", func() {
		command := exec.Command(pathToOmenCLI, "-u="+validUsername)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Eventually(session, opsManagerTimeOut).Should(gexec.Exit(1))
	})

	It("should print an error and exit with nonzero value if command parameter is not recognised", func() {
		command := exec.Command(pathToOmenCLI, "-u="+validUsername, "thisIsNotCommand")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Eventually(session, opsManagerTimeOut).Should(gexec.Exit(1))
	})
})
