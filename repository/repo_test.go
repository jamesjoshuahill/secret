package repository_test

import (
	"github.com/jamesjoshuahill/ciphers/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	It("stores and retrieves ciphers", func() {
		repo := repository.New()

		cipher := repository.Cipher{
			ID:         "some-id",
			CipherText: "some-cipher-text",
		}
		err := repo.Store(cipher)
		Expect(err).NotTo(HaveOccurred())

		actualCipher, err := repo.FindByID("some-id")
		Expect(err).NotTo(HaveOccurred())

		Expect(actualCipher).To(Equal(cipher))
	})

	It("fails when the cipher already exits", func() {
		repo := repository.New()

		err := repo.Store(repository.Cipher{
			ID:         "some-id",
			CipherText: "some-cipher-text",
		})
		Expect(err).NotTo(HaveOccurred())

		err = repo.Store(repository.Cipher{
			ID:         "some-id",
			CipherText: "some-cipher-text",
		})
		Expect(err).To(MatchError("already exists"))
	})

	It("fails when it cannot find a cipher", func() {
		repo := repository.New()

		_, err := repo.FindByID("some-id")
		Expect(err).To(MatchError("not found"))
	})
})
