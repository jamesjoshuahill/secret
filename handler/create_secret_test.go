package handler_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jamesjoshuahill/secret/aes"
	"github.com/jamesjoshuahill/secret/handler"
	"github.com/jamesjoshuahill/secret/handler/fake"
	"github.com/jamesjoshuahill/secret/inmemory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateSecret", func() {
	var (
		repo    *fake.Repo
		encrypt handler.EncryptFunc
		res     *httptest.ResponseRecorder
		req     *http.Request
	)

	BeforeEach(func() {
		repo = new(fake.Repo)
		res = httptest.NewRecorder()

		var err error
		req, err = http.NewRequest("POST", "/v1/secrets", strings.NewReader(`{
			"id": "client-secret-id",
			"data": "some plain text"
		}`))
		Expect(err).NotTo(HaveOccurred())
		req.Header.Set("Content-Type", "application/json")
	})

	It("encrypts the plain text", func() {
		var plainTextReceived string
		encrypt = func(plainText string) (aes.Secret, error) {
			plainTextReceived = plainText
			return aes.Secret{}, nil
		}
		handler := handler.CreateSecret{Repository: repo, Encrypt: encrypt}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(plainTextReceived).To(Equal("some plain text"))
	})

	It("fails when the request content type is not JSON", func() {
		req.Header.Set("Content-Type", "text/plain")
		handler := handler.CreateSecret{Repository: repo, Encrypt: encrypt}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusUnsupportedMediaType))
		Expect(res.Body.String()).To(ContainSubstring("unsupported Content-Type"))
	})

	It("fails when the request body cannot be parsed", func() {
		req.Body = ioutil.NopCloser(strings.NewReader("not json"))
		handler := handler.CreateSecret{Repository: repo, Encrypt: encrypt}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusBadRequest))
		Expect(res.Body.String()).To(ContainSubstring("decoding request body"))
	})

	It("fails when the plain text cannot be encrypted", func() {
		encrypt = func(string) (aes.Secret, error) {
			return aes.Secret{}, errors.New("fake error")
		}
		handler := handler.CreateSecret{Repository: repo, Encrypt: encrypt}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusInternalServerError), res.Body.String())
		Expect(res.Body.String()).To(ContainSubstring("encrypting data"))
	})

	It("stores the secret", func() {
		encrypt = func(string) (aes.Secret, error) {
			return aes.Secret{
				Key:        "key for client-secret-id",
				Nonce:      "some nonce",
				CipherText: "some cipher text",
			}, nil
		}
		handler := handler.CreateSecret{Repository: repo, Encrypt: encrypt}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(repo.StoreCall.Received.Secret).To(Equal(inmemory.Secret{
			ID:         "client-secret-id",
			Nonce:      "some nonce",
			CipherText: "some cipher text",
		}))
	})

	It("fails when the secret already exists", func() {
		repo.StoreCall.Returns.Error = errors.New("fake error")
		handler := handler.CreateSecret{Repository: repo, Encrypt: encrypt}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusConflict), res.Body.String())
		Expect(res.Body.String()).To(ContainSubstring("secret already exists"))
	})
})
