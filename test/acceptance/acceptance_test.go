package acceptance

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("omen acceptance", func() {

	const timeout = 0.1

	It("should print an error if the url argument is missing", func() {
		command := exec.Command(pathToOmenCLI, "-u=user", "apply-changes")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, timeout).Should(gexec.Exit(1))
		Eventually(func() string { return string(session.Out.Contents()) }, timeout).Should(ContainSubstring("Opsman host is required"))
	})

	Describe("usernames and passwords", func() {
		It("should require both to be present", func() {
			command := exec.Command(pathToOmenCLI, "-t=https://127.0.0.1", "-u=user", "apply-changes")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, timeout).Should(gexec.Exit(1))
			Eventually(func() string { return string(session.Out.Contents()) }, timeout).Should(ContainSubstring("Opsman user secret is required"))
		})
	})

	Describe("client ids and secrets", func() {
		It("should require both to be present", func() {
			command := exec.Command(pathToOmenCLI, "-t=https://127.0.0.1", "-c=client", "apply-changes")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, timeout).Should(gexec.Exit(1))
			Eventually(func() string { return string(session.Out.Contents()) }, timeout).Should(ContainSubstring("Opsman client secret is required."))
		})
	})
})
