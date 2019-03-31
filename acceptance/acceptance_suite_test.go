package acceptance_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const serverPort = 8080

var (
	pathToServerBinary string
	client             *http.Client
)

func TestAcceptance(t *testing.T) {
	var serverSession *gexec.Session

	BeforeSuite(func() {
		var err error
		pathToServerBinary, err = gexec.Build("github.com/jamesjoshuahill/ciphers/cmd/server")
		Expect(err).NotTo(HaveOccurred())

		serverSession = startServer(pathToServerBinary)

		Eventually(dialServer).Should(Succeed())

		client = newClient()
	})

	AfterSuite(func() {
		serverSession.Terminate().Wait()

		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func startServer(pathToServerBinary string) *gexec.Session {
	cmd := exec.Command(
		pathToServerBinary,
		fmt.Sprintf("--port=%d", serverPort),
		fmt.Sprintf("--cert=fixtures/cert.pem"),
		fmt.Sprintf("--key=fixtures/key.pem"),
	)

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}

func dialServer() error {
	_, err := net.Dial("tcp", fmt.Sprintf(":%d", serverPort))
	return err
}

func serverUrl(path string) string {
	return fmt.Sprintf("https://127.0.0.1:%d/%s", serverPort, path)
}

func newClient() *http.Client {
	certPool := x509.NewCertPool()

	rootCA, err := ioutil.ReadFile("fixtures/cert.pem")
	Expect(err).NotTo(HaveOccurred())

	ok := certPool.AppendCertsFromPEM(rootCA)
	Expect(ok).To(BeTrue(), "failed to append root CA cert")

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
}
