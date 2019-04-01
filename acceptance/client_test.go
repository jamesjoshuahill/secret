package acceptance_test

import (
	"github.com/jamesjoshuahill/ciphers/pkg/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	It("stores and retrieves ciphers", func() {
		c := client.New(serverBaseURL(), httpsClient)

		By("storing a cipher")
		id := []byte("my-id")
		payload := []byte("my-payload")
		key, err := c.Store(id, payload)

		Expect(err).NotTo(HaveOccurred())
		Expect(key).NotTo(BeEmpty())

		By("retreiving a cipher")
		actualPayload, err := c.Retrieve(id, key)

		Expect(err).NotTo(HaveOccurred())
		Expect(actualPayload).To(Equal(payload))
	})
})
