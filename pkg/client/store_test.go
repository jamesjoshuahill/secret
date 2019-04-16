package client_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/xerrors"

	"github.com/jamesjoshuahill/ciphers/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {
	It("makes valid create secret requests", func() {
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
			ContainSubstring("create secret request"),
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
		unerr := &client.UnexpectedError{}
		Expect(xerrors.As(err, unerr)).To(BeTrue())
		Expect(unerr.StatusCode).To(Equal(http.StatusInternalServerError))
		Expect(unerr.Message).To(Equal("fake error"))
	})

	It("fails when the response is not unexpected and malformed", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       readCloser("not json"),
		}

		_, err := c.Store([]byte("some-id"), []byte("some-payload"))

		Expect(err).To(HaveOccurred())
		unerr := &client.UnexpectedError{}
		Expect(xerrors.As(err, unerr)).To(BeTrue())
		Expect(unerr.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("fails when the secret already exists", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusConflict,
			Body:       readCloser(`{"error":"secret already exists"}`),
		}

		_, err := c.Store([]byte("some-id"), []byte("some-payload"))

		Expect(err).To(HaveOccurred())
		unerr := err.(alreadyExists)
		Expect(unerr.AlreadyExists()).To(BeTrue())
	})

	It("fails when the response cannot be parsed", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusOK,
			Body:       readCloser("not json"),
		}

		_, err := c.Store([]byte("some-id"), []byte("some-payload"))

		Expect(err).To(MatchError(
			ContainSubstring("decoding create secret response body"),
		))
	})
})
