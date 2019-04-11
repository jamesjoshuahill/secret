package inmemory_test

import (
	"github.com/jamesjoshuahill/ciphers/repository"
	"github.com/jamesjoshuahill/ciphers/repository/inmemory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	It("stores and retrieves ciphers", func() {
		repo := inmemory.New()

		cipher := repository.Cipher{
			ID:         "some-id",
			Nonce:      "some-nonce",
			CipherText: "some-cipher-text",
		}
		err := repo.Store(cipher)
		Expect(err).NotTo(HaveOccurred())

		actualCipher, err := repo.FindByID("some-id")
		Expect(err).NotTo(HaveOccurred())

		Expect(actualCipher).To(Equal(cipher))
	})

	It("fails when the cipher already exits", func() {
		repo := inmemory.New()

		err := repo.Store(repository.Cipher{ID: "some-id"})
		Expect(err).NotTo(HaveOccurred())

		err = repo.Store(repository.Cipher{ID: "some-id"})
		Expect(err).To(MatchError("already exists"))
	})

	It("fails when it cannot find a cipher", func() {
		repo := inmemory.New()

		_, err := repo.FindByID("some-id")
		Expect(err).To(MatchError("not found"))
	})
})
