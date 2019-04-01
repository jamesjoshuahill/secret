package acceptance_test

import (
	"github.com/jamesjoshuahill/ciphers/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	It("stores ciphers", func() {
		var serverClient client.ServerClient
		serverClient = httpsClient
		c := client.New(serverBaseURL(), serverClient)

		key, err := c.Store([]byte("my-id"), []byte("my-payload"))

		Expect(err).NotTo(HaveOccurred())
		Expect(key).NotTo(BeEmpty())
	})
})
