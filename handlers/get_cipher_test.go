package handlers_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"

	"github.com/jamesjoshuahill/ciphers/repository"

	"github.com/jamesjoshuahill/ciphers/handlers/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jamesjoshuahill/ciphers/handlers"
)

var _ = Describe("GetCipher", func() {
	var (
		repo      *fakes.FakeRepo
		decrypter *fakes.FakeDecrypter
		res       *httptest.ResponseRecorder
		req       *http.Request
		router    *mux.Router
	)

	BeforeEach(func() {
		repo = new(fakes.FakeRepo)
		decrypter = new(fakes.FakeDecrypter)
		res = httptest.NewRecorder()
		router = mux.NewRouter()

		var err error
		req, err = http.NewRequest("GET", "/v1/ciphers/client-cipher-id", strings.NewReader(`{
			"key": "key for client-cipher-id"
		}`))
		Expect(err).NotTo(HaveOccurred())
		req.Header.Set("Content-Type", "application/json")
	})

	It("retrieves the cipher", func() {
		repo.FindByResourceIDCall.Returns.Cipher = repository.Cipher{
			ID:         "client-cipher-id",
			CipherText: "some cipher text",
		}
		router.Handle("/v1/ciphers/{id}", &handlers.GetCipher{Repository: repo, Decrypter: decrypter})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(repo.FindByResourceIDCall.Received.ID).To(Equal("client-cipher-id"))
	})

	It("decrypts the ciphertext", func() {
		repo.FindByResourceIDCall.Returns.Cipher = repository.Cipher{
			ID:         "client-cipher-id",
			CipherText: "some cipher text",
		}
		router.Handle("/v1/ciphers/{id}", &handlers.GetCipher{Repository: repo, Decrypter: decrypter})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(decrypter.DecryptCall.Received.Key).To(Equal("key for client-cipher-id"))
		Expect(decrypter.DecryptCall.Received.CipherText).To(Equal("some cipher text"))
	})

	It("fails when the request content type is not JSON", func() {
		req.Header.Set("Content-Type", "text/plain")
		handler := handlers.GetCipher{Repository: repo, Decrypter: decrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusUnsupportedMediaType))
		Expect(res.Body.String()).To(ContainSubstring("unsupported Content-Type"))
	})

	It("fails when the request body cannot be parsed", func() {
		req.Body = ioutil.NopCloser(strings.NewReader("not json"))
		router.Handle("/v1/ciphers/{id}", &handlers.GetCipher{Repository: repo, Decrypter: decrypter})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusBadRequest), res.Body.String())
		Expect(res.Body.String()).To(ContainSubstring("decoding request body"))
	})

	It("fails when the cipher is not found", func() {
		repo.FindByResourceIDCall.Returns.Error = errors.New("fake error")
		router.Handle("/v1/ciphers/{id}", &handlers.GetCipher{Repository: repo, Decrypter: decrypter})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusNotFound), res.Body.String())
		Expect(res.Body.String()).To(ContainSubstring("not found"))
	})

	It("fails when the cipher cannot be decrypted", func() {
		decrypter.DecryptCall.Returns.Error = errors.New("fake error")
		handler := handlers.GetCipher{Repository: repo, Decrypter: decrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusUnauthorized))
		Expect(res.Body.String()).To(ContainSubstring("wrong key"))
	})
})
