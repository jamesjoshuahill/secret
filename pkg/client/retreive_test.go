package client_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jamesjoshuahill/secret/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Retrieve", func() {
	It("makes valid get secret requests", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusOK,
			Body:       readCloser(`{"data":"some-payload"}`),
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
			ContainSubstring("get secret request"),
			ContainSubstring("fake error"),
		)))
	})

	It("fails when response is not unexpected", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       readCloser(`{"error":"fake error"}`),
		}

		_, err := c.Retrieve([]byte("some-id"), []byte("some-key"))

		Expect(err).To(HaveOccurred())
		unerr := &client.UnexpectedError{}
		Expect(errors.As(err, unerr)).To(BeTrue())
		Expect(unerr.StatusCode).To(Equal(http.StatusInternalServerError))
		Expect(unerr.Message).To(Equal("fake error"))
	})

	It("fails when the response is not unexpected and malformed", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       readCloser("not json"),
		}

		_, err := c.Retrieve([]byte("some-id"), []byte("some-key"))

		Expect(err).To(HaveOccurred())
		unerr := &client.UnexpectedError{}
		Expect(errors.As(err, unerr)).To(BeTrue())
		Expect(unerr.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("fails with the wrong secret id or key", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Body:       readCloser(`{"error":"wrong id or key"}`),
		}

		_, err := c.Retrieve([]byte("some-id"), []byte("some-key"))

		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, client.ErrWrongIDOrKey)).To(BeTrue())
	})

	It("fails when the response cannot be parsed", func() {
		httpsClient.DoCall.Returns.Response = &http.Response{
			StatusCode: http.StatusOK,
			Body:       readCloser("not json"),
		}

		_, err := c.Retrieve([]byte("some-id"), []byte("some-key"))

		Expect(err).To(MatchError(
			ContainSubstring("decoding get secret response body"),
		))
	})
})
