package client_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/pkg/client"
	"github.com/jamesjoshuahill/ciphers/pkg/client/fakes"
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

type unexpectedError interface {
	StatusCode() int
	Message() string
}

var _ = Describe("Client", func() {
	const baseURL = "https://example.com:8080"

	var (
		c           Client
		httpsClient *fakes.FakeHTTPSClient
	)

	BeforeEach(func() {
		httpsClient = new(fakes.FakeHTTPSClient)
		c = client.New(baseURL, httpsClient)
	})

	Context("Store", func() {
		It("makes valid create cipher requests", func() {
			httpsClient.DoCall.Returns.Response = &http.Response{
				StatusCode: http.StatusOK,
				Body:       readCloser(`{"key":"some-key"}`),
			}

			key, err := c.Store([]byte("some-id"), []byte("some-payload"))

			Expect(err).NotTo(HaveOccurred())
			req := httpsClient.DoCall.Received.Request
			Expect(req.Method).To(Equal("POST"))
			Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
			body, err := ioutil.ReadAll(req.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(MatchJSON(`{
				"id": "some-id",
				"data": "some-payload"
			}`))
			Expect(key).To(Equal([]byte("some-key")))
		})

		It("fails when the request errors", func() {
			httpsClient.DoCall.Returns.Error = errors.New("fake error")

			_, err := c.Store([]byte("some-id"), []byte("some-payload"))

			Expect(err).To(MatchError(SatisfyAll(
				ContainSubstring("create cipher request"),
				ContainSubstring("fake error"),
			)))
		})

		It("fails when response is not unexpected", func() {
			httpsClient.DoCall.Returns.Response = &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       readCloser(`{"error":"fake error"}`),
			}

			_, err := c.Store([]byte("some-id"), []byte("some-payload"))

			Expect(err).To(HaveOccurred())
			unerr := err.(unexpectedError)
			Expect(unerr.StatusCode()).To(Equal(http.StatusInternalServerError))
			Expect(unerr.Message()).To(Equal("fake error"))
		})

		It("fails when the response is not unexpected and malformed", func() {
			httpsClient.DoCall.Returns.Response = &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       readCloser("not json"),
			}

			_, err := c.Store([]byte("some-id"), []byte("some-payload"))

			Expect(err).To(HaveOccurred())
			unerr := err.(unexpectedError)
			Expect(unerr.StatusCode()).To(Equal(http.StatusInternalServerError))
		})

		It("fails when the response cannot be parsed", func() {
			httpsClient.DoCall.Returns.Response = &http.Response{
				StatusCode: http.StatusOK,
				Body:       readCloser("not json"),
			}

			_, err := c.Store([]byte("some-id"), []byte("some-payload"))

			Expect(err).To(MatchError(
				ContainSubstring("decoding create cipher response body"),
			))
		})
	})

	Context("Retrieve", func() {
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
})
