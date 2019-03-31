package handlers_test

import (
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
	It("retrieves the cipher", func() {
		repo := new(fakes.FakeRepo)
		repo.FindByResourceIDCall.Returns.Cipher = repository.Cipher{
			ResourceID: "client-cipher-id",
			Data:       "some plain text",
			Key:        "key for client-cipher-id",
		}
		req, err := http.NewRequest("GET", "/v1/ciphers/client-cipher-id", strings.NewReader(`{
			"key": "key for client-cipher-id"
		}`))
		Expect(err).NotTo(HaveOccurred())
		res := httptest.NewRecorder()
		router := mux.NewRouter()
		router.Handle("/v1/ciphers/{resource_id}", &handlers.GetCipher{Repository: repo})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(repo.FindByResourceIDCall.Received.ResourceID).To(Equal("client-cipher-id"))
	})
})
