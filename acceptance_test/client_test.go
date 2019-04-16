package acceptance_test

import (
	"github.com/jamesjoshuahill/ciphers/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	It("stores and retrieves secrets", func() {
		c := client.NewClient(serverBaseURL(), httpsClient)

		By("storing a secret")
		id := []byte("my-id")
		payload := []byte("my-payload")
		key, err := c.Store(id, payload)

		Expect(err).NotTo(HaveOccurred())
		Expect(key).NotTo(BeEmpty())

		By("retrieving a secret")
		actualPayload, err := c.Retrieve(id, key)

		Expect(err).NotTo(HaveOccurred())
		Expect(actualPayload).To(Equal(payload))
	})
})
