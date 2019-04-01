package client_test

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/jamesjoshuahill/ciphers/pkg/client"
	"github.com/jamesjoshuahill/ciphers/pkg/client/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const baseURL = "https://example.com:8080"

var (
	c           Client
	httpsClient *fakes.FakeHTTPSClient
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

type unexpectedError interface {
	StatusCode() int
	Message() string
}

type alreadyExistsError interface {
	AlreadyExists() bool
}

func TestClient(t *testing.T) {
	BeforeEach(func() {
		httpsClient = new(fakes.FakeHTTPSClient)
		c = client.New(baseURL, httpsClient)
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

func readCloser(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}
