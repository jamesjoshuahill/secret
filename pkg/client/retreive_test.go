package client_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Retrieve", func() {
	It("makes valid get cipher requests", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			Body: readCloser(`{"data":"some-payload"}`),
		}

		actualPayload, err := c.Retrieve([]byte("some-id"), []byte("some-key"))

		Expect(err).NotTo(HaveOccurred())
		req := httpsClient.DoCall.Received.Request
		Expect(req.Method).To(Equal("GET"))
		Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
		body, err := ioutil.ReadAll(req.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(MatchJSON(`{
				"key": "some-key"
			}`))
		Expect(actualPayload).To(Equal([]byte("some-payload")))
	})

	It("fails when the request errors", func() {
		httpsClient.DoCall.Returns.Error = errors.New("fake error")

		_, err := c.Retrieve([]byte("some-id"), []byte("some-key"))

		Expect(err).To(MatchError(SatisfyAll(
			ContainSubstring("get cipher request"),
			ContainSubstring("fake error"),
		)))
	})

	It("fails when the response cannot be parsed", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			Body: readCloser("not json"),
		}

		_, err := c.Retrieve([]byte("some-id"), []byte("some-key"))

		Expect(err).To(MatchError(
			ContainSubstring("decoding get cipher response body"),
		))
	})
})
