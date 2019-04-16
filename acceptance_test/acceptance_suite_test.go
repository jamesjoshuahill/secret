package acceptance_test

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"testing"

	"github.com/jamesjoshuahill/ciphers/pkg/client"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const serverPort = 8080

var (
	pathToServerBinary string
	httpsClient        *http.Client
)

func TestAcceptance(t *testing.T) {
	var serverSession *gexec.Session

	BeforeSuite(func() {
		var err error
		pathToServerBinary, err = gexec.Build("github.com/jamesjoshuahill/ciphers/cmd/secret-server")
		Expect(err).NotTo(HaveOccurred())

		serverSession = startServer(pathToServerBinary)

		Eventually(dialServer).Should(Succeed())

		httpsClient = newClient()
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
		fmt.Sprintf("--cert=testdata/cert.pem"),
		fmt.Sprintf("--key=testdata/key.pem"),
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
	return fmt.Sprintf("%s/%s", serverBaseURL(), path)
}

func serverBaseURL() string {
	return fmt.Sprintf("https://127.0.0.1:%d", serverPort)
}

func newClient() *http.Client {
	certPool := x509.NewCertPool()

	rootCA, err := ioutil.ReadFile("testdata/cert.pem")
	Expect(err).NotTo(HaveOccurred())

	ok := certPool.AppendCertsFromPEM(rootCA)
	Expect(ok).To(BeTrue(), "failed to append root CA cert")

	return client.DefaultHTTPSClient(certPool)
}
