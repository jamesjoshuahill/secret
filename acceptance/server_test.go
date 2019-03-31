package acceptance_test

import (
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	It("accepts a valid create cipher request", func() {
		res, err := http.Post(serverUrl("v1/ciphers"), "application/json", strings.NewReader(`{
			"data": "some plain text",
			"resource_id": "client cipher id"
		}`))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(MatchJSON(`{
			"id": "server cipher id",
			"data": "some plain text",
			"resource_id": "client cipher id"
		}`))
	})

	It("rejects a malformed create cipher request", func() {
		res, err := http.Post(serverUrl("v1/ciphers"), "application/json", strings.NewReader("not json"))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(ContainSubstring("error decoding request body"))
	})
})
