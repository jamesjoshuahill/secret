package client_test

import (
	"github.com/jamesjoshuahill/ciphers/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var _ = Describe("Client", func() {
	It("is a client", func() {
		var c Client
		c = client.New("", nil)
		Expect(c).NotTo(BeNil())
	})
})
