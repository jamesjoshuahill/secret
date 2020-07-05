package client_test

import (
	"io"
	"strings"
	"testing"

	"github.com/jamesjoshuahill/secret/pkg/client"
	"github.com/jamesjoshuahill/secret/pkg/client/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const baseURL = "https://example.com:8080"

var (
	c           Client
	httpsClient *fake.HTTPSClient
)

// Client provides functionality to interact with the encryption-server
type Client interface {
	// Store accepts an id and a payload in bytes and requests that the
	// encryption-server stores them in its data store
	Store(id, payload []byte) (aesKey []byte, err error)

	// Retrieve accepts an id and an AES key, and requests that the
	// encryption-server retrieves the original (decrypted) bytes stored
	// with the provided id
	Retrieve(id, aesKey []byte) (payload []byte, err error)
}

func TestClient(t *testing.T) {
	BeforeEach(func() {
		httpsClient = new(fake.HTTPSClient)
		c = client.NewClient(baseURL, httpsClient)
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

type ReadCloserSpy struct {
	io.Reader
	CloseCount int
}

func (r *ReadCloserSpy) Close() error {
	r.CloseCount++
	return nil
}

func NewReadCloserSpy(s string) *ReadCloserSpy {
	return &ReadCloserSpy{Reader: strings.NewReader(s)}
}
