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
			"resource_id": "client-cipher-id",
			"data": "some plain text"
		}`))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(MatchJSON(`{
			"resource_id": "client-cipher-id",
			"key": "key for server-cipher-id"
		}`))
	})

	It("rejects a malformed create cipher request", func() {
		res, err := http.Post(serverUrl("v1/ciphers"), "application/json", strings.NewReader("not json"))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(SatisfyAll(
			ContainSubstring("error"),
			ContainSubstring("decoding request body"),
		))
	})

	It("accepts a valid get cipher request", func() {
		req, err := http.NewRequest("GET", serverUrl("v1/ciphers/client-cipher-id"), strings.NewReader(`{
			"key": "key for server-cipher-id"
		}`))
		Expect(err).NotTo(HaveOccurred())

		res, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(MatchJSON(`{
			"resource_id": "client-cipher-id",
			"data": "some plain text"
		}`))
	})

	It("rejects a malformed get cipher request", func() {
		req, err := http.NewRequest("GET", serverUrl("v1/ciphers/client-cipher-id"), strings.NewReader("not json"))
		Expect(err).NotTo(HaveOccurred())

		res, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(SatisfyAll(
			ContainSubstring("error"),
			ContainSubstring("decoding request body"),
		))
	})

	It("rejects a get cipher request with the wrong key", func() {
		req, err := http.NewRequest("GET", serverUrl("v1/ciphers/client-cipher-id"), strings.NewReader(`{
			"key": "wrong key"
		}`))
		Expect(err).NotTo(HaveOccurred())

		res, err := http.DefaultClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusUnauthorized))
		Expect(res.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(SatisfyAll(
			ContainSubstring("error"),
			ContainSubstring("wrong key"),
		))
	})
})
