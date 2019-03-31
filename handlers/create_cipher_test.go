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

		req, err := http.NewRequest("POST", "/v1/ciphers", strings.NewReader(`{
			"resource_id": "client-cipher-id",
			"data": "some plain text"
		}`))
		Expect(err).NotTo(HaveOccurred())

		res := httptest.NewRecorder()
		handler := handlers.CreateCipher{Repository: repo}
		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(repo.StoreCall.Received.Cipher).To(Equal(repository.Cipher{
			ResourceID: "client-cipher-id",
			Data:       "some plain text",
			Key:        "key for client-cipher-id",
		}))
	})
})
