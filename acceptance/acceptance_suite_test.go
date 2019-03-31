package acceptance_test

import (
	"fmt"
	"net"
	"os/exec"
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const serverPort = 8080

func TestAcceptance(t *testing.T) {
	var (
		pathToServerBinary string
		serverSession      *gexec.Session
	)

	BeforeSuite(func() {
		var err error
		pathToServerBinary, err = gexec.Build("github.com/jamesjoshuahill/encryption/cmd/server")
		Expect(err).NotTo(HaveOccurred())

		serverSession = startServer(pathToServerBinary)
	})

	AfterSuite(func() {
		serverSession.Terminate().Wait()

		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func startServer(pathToServerBinary string) *gexec.Session {
	cmd := exec.Command(pathToServerBinary)
	cmd.Env = []string{fmt.Sprintf("PORT=%d", serverPort)}

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}

func dialServer() error {
	_, err := net.Dial("tcp", fmt.Sprintf(":%d", serverPort))
	return err
}
