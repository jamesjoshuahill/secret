package client_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jamesjoshuahill/secret/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {
	It("makes valid create secret requests", func() {
		resBody := NewReadCloserSpy(`{"key":"some-key"}`)
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusOK,
			Body:       resBody,
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
		Expect(resBody.CloseCount).To(Equal(1))
	})

	It("fails when the request errors", func() {
		httpsClient.DoCall.Returns.Error = errors.New("fake error")

		_, err := c.Store([]byte("some-id"), []byte("some-payload"))

		Expect(err).To(MatchError(SatisfyAll(
			ContainSubstring("create secret request"),
			ContainSubstring("fake error"),
		)))
	})

	It("fails when response is unexpected", func() {
		body := NewReadCloserSpy(`{"error":"fake error"}`)
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       body,
		}

		_, err := c.Store([]byte("some-id"), []byte("some-payload"))

		Expect(err).To(HaveOccurred())
		unerr := &client.UnexpectedError{}
		Expect(errors.As(err, unerr)).To(BeTrue())
		Expect(unerr.StatusCode).To(Equal(http.StatusInternalServerError))
		Expect(unerr.Message).To(Equal("fake error"))
		Expect(body.CloseCount).To(Equal(1))
	})

	It("fails when the response is unexpected and malformed", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       NewReadCloserSpy("not json"),
		}

		_, err := c.Store([]byte("some-id"), []byte("some-payload"))

		Expect(err).To(HaveOccurred())
		unerr := &client.UnexpectedError{}
		Expect(errors.As(err, unerr)).To(BeTrue())
		Expect(unerr.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("fails when the secret already exists", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusConflict,
			Body:       NewReadCloserSpy(`{"error":"secret already exists"}`),
		}

		_, err := c.Store([]byte("some-id"), []byte("some-payload"))

		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, client.ErrAlreadyExists)).To(BeTrue())
	})

	It("fails when the response cannot be parsed", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusOK,
			Body:       NewReadCloserSpy("not json"),
		}

		_, err := c.Store([]byte("some-id"), []byte("some-payload"))

		Expect(err).To(MatchError(
			ContainSubstring("decoding create secret response body"),
		))
	})
})
