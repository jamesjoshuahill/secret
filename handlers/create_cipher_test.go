package handlers_test

import (
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
	It("stores the cipher", func() {
		repo := new(fakes.FakeRepo)

		req, err := http.NewRequest("POST", "v1/ciphers", strings.NewReader(`{
			"resource_id": "client-cipher-id",
			"data": "some plaintext"
		}`))
		Expect(err).NotTo(HaveOccurred())

		recorder := httptest.NewRecorder()
		handler := handlers.CreateCipher{Repository: repo}
		handler.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusOK), recorder.Body.String())
		Expect(repo.StoreCall.Received.Cipher).To(Equal(repository.Cipher{
			ResourceID: "client-cipher-id",
			Data:       "some plaintext",
			Key:        "key for server-cipher-id",
		}))
	})
})
