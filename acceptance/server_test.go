package acceptance_test

import (
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	It("accepts an encryption request", func() {
		res, err := http.Post(serverUrl("v1/ciphers"), "application/json", strings.NewReader(`{
			"data": "some plain text",
			"resource_id": "client cipher id"
		}`))
		Expect(err).NotTo(HaveOccurred())

		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})
})
