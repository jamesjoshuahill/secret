package handlers_test

import (
	"errors"
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
		repo   *fakes.FakeRepo
		res    *httptest.ResponseRecorder
		req    *http.Request
		router *mux.Router
	)

	BeforeEach(func() {
		repo = new(fakes.FakeRepo)
		res = httptest.NewRecorder()
		router = mux.NewRouter()

		var err error
		req, err = http.NewRequest("GET", "/v1/ciphers/client-cipher-id", strings.NewReader(`{
			"key": "key for client-cipher-id"
		}`))
		Expect(err).NotTo(HaveOccurred())
	})

	It("retrieves the cipher", func() {
		repo.FindByResourceIDCall.Returns.Cipher = repository.Cipher{
			ResourceID: "client-cipher-id",
			CipherText: "some plain text",
			Key:        "key for client-cipher-id",
		}
		router.Handle("/v1/ciphers/{resource_id}", &handlers.GetCipher{Repository: repo})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(repo.FindByResourceIDCall.Received.ResourceID).To(Equal("client-cipher-id"))
	})

	It("fails when the cipher cannot be retrieved", func() {
		repo.FindByResourceIDCall.Returns.Error = errors.New("fake error")
		router.Handle("/v1/ciphers/{resource_id}", &handlers.GetCipher{Repository: repo})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusInternalServerError), res.Body.String())
		Expect(res.Body.String()).To(ContainSubstring("finding cipher"))
	})
})
