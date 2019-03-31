package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jamesjoshuahill/ciphers/handlers/fakes"

	"github.com/jamesjoshuahill/ciphers/repository"

	"github.com/jamesjoshuahill/ciphers/handlers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateCipher", func() {
	var (
		repo      *fakes.FakeRepo
		encrypter *fakes.FakeEncrypter
		res       *httptest.ResponseRecorder
		req       *http.Request
	)

	BeforeEach(func() {
		repo = new(fakes.FakeRepo)
		encrypter = new(fakes.FakeEncrypter)
		res = httptest.NewRecorder()

		var err error
		req, err = http.NewRequest("POST", "/v1/ciphers", strings.NewReader(`{
			"resource_id": "client-cipher-id",
			"data": "some plain text"
		}`))
		Expect(err).NotTo(HaveOccurred())
	})

	It("encrypts the plain text", func() {
		encrypter.EncryptCall.Returns.Key = "key for client-cipher-id"
		handler := handlers.CreateCipher{Repository: repo, Encrypter: encrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(encrypter.EncryptCall.Received.Plaintext).To(Equal("some plain text"))
	})

	It("stores the cipher", func() {
		encrypter.EncryptCall.Returns.Key = "key for client-cipher-id"
		handler := handlers.CreateCipher{Repository: repo, Encrypter: encrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(repo.StoreCall.Received.Cipher).To(Equal(repository.Cipher{
			ResourceID: "client-cipher-id",
			Data:       "some plain text",
			Key:        "key for client-cipher-id",
		}))
	})

	It("fails when the cipher cannot be stored", func() {
		encrypter.EncryptCall.Returns.Key = "key for client-cipher-id"
		repo.StoreCall.Returns.Error = errors.New("fake error")
		handler := handlers.CreateCipher{Repository: repo, Encrypter: encrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusInternalServerError), res.Body.String())
		Expect(res.Body.String()).To(ContainSubstring("storing cipher"))
	})
})
