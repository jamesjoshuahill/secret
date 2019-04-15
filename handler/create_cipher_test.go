package handler_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jamesjoshuahill/ciphers/repository"

	"github.com/jamesjoshuahill/ciphers/encryption"

	"github.com/jamesjoshuahill/ciphers/handler/fake"

	"github.com/jamesjoshuahill/ciphers/handler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateCipher", func() {
	var (
		repo      *fake.FakeRepo
		encrypter *fake.FakeEncrypter
		res       *httptest.ResponseRecorder
		req       *http.Request
	)

	BeforeEach(func() {
		repo = new(fake.FakeRepo)
		encrypter = new(fake.FakeEncrypter)
		res = httptest.NewRecorder()

		var err error
		req, err = http.NewRequest("POST", "/v1/ciphers", strings.NewReader(`{
			"id": "client-cipher-id",
			"data": "some plain text"
		}`))
		Expect(err).NotTo(HaveOccurred())
		req.Header.Set("Content-Type", "application/json")
	})

	It("encrypts the plain text", func() {
		handler := handler.CreateCipher{Repository: repo, Encrypter: encrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(encrypter.EncryptCall.Received.PlainText).To(Equal("some plain text"))
	})

	It("fails when the request content type is not JSON", func() {
		req.Header.Set("Content-Type", "text/plain")
		handler := handler.CreateCipher{Repository: repo, Encrypter: encrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusUnsupportedMediaType))
		Expect(res.Body.String()).To(ContainSubstring("unsupported Content-Type"))
	})

	It("fails when the request body cannot be parsed", func() {
		req.Body = ioutil.NopCloser(strings.NewReader("not json"))
		handler := handler.CreateCipher{Repository: repo, Encrypter: encrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusBadRequest))
		Expect(res.Body.String()).To(ContainSubstring("decoding request body"))
	})

	It("fails when the plain text cannot be encrypted", func() {
		encrypter.EncryptCall.Returns.Error = errors.New("fake error")
		handler := handler.CreateCipher{Repository: repo, Encrypter: encrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusInternalServerError), res.Body.String())
		Expect(res.Body.String()).To(ContainSubstring("encrypting data"))
	})

	It("stores the cipher", func() {
		encrypter.EncryptCall.Returns.Cipher = encryption.Cipher{
			Key:        "key for client-cipher-id",
			Nonce:      "some nonce",
			CipherText: "some cipher text",
		}
		handler := handler.CreateCipher{Repository: repo, Encrypter: encrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(repo.StoreCall.Received.Cipher).To(Equal(repository.Cipher{
			ID:         "client-cipher-id",
			Nonce:      "some nonce",
			CipherText: "some cipher text",
		}))
	})

	It("fails when the cipher already exists", func() {
		repo.StoreCall.Returns.Error = errors.New("fake error")
		handler := handler.CreateCipher{Repository: repo, Encrypter: encrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusConflict), res.Body.String())
		Expect(res.Body.String()).To(ContainSubstring("cipher already exists"))
	})
})
