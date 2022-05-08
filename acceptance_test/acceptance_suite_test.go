package acceptance_test

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"testing"
	"time"

	"github.com/jamesjoshuahill/secret/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const (
	serverHost = "localhost"
	serverPort = 8080
)

var (
	pathToServerBinary string
	httpsClient        *http.Client
)

func TestAcceptance(t *testing.T) {
	var serverSession *gexec.Session

	BeforeSuite(func() {
		var err error
		pathToServerBinary, err = gexec.Build("github.com/jamesjoshuahill/secret/cmd/secret-server")
		Expect(err).NotTo(HaveOccurred())

		serverSession = startServer(pathToServerBinary)

		Eventually(dialServer, 5*time.Second).Should(Succeed())

		httpsClient = newClient()
	})

	AfterSuite(func() {
		serverSession.Terminate().Wait(time.Second * (6 + 1))

		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func startServer(pathToServerBinary string) *gexec.Session {
	cmd := exec.Command(
		pathToServerBinary,
		"--host", "localhost",
		"--port", fmt.Sprintf("%d", serverPort),
		"--cert", "testdata/cert.pem",
		"--key", "testdata/key.pem",
	)

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}

func dialServer() error {
	_, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverHost, serverPort))
	return err
}

func serverURL(path string) string {
	return fmt.Sprintf("%s/%s", serverBaseURL(), path)
}

func serverBaseURL() string {
	return fmt.Sprintf("https://%s:%d", serverHost, serverPort)
}

func newClient() *http.Client {
	certPool := x509.NewCertPool()

	rootCA, err := ioutil.ReadFile("testdata/cert.pem")
	Expect(err).NotTo(HaveOccurred())

	ok := certPool.AppendCertsFromPEM(rootCA)
	Expect(ok).To(BeTrue(), "failed to append root CA cert")

	return client.DefaultHTTPSClient(certPool)
}
