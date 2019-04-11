package inmemory_test

import (
	"sync"

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

	It("does not store a cipher more than once", func() {
		repo := inmemory.New()

		cipher := repository.Cipher{
			ID:         "some-id",
			Nonce:      "some-nonce",
			CipherText: "some-cipher-text",
		}

		wg := sync.WaitGroup{}
		wg.Add(3)

		var err1, err2, err3 error
		go func() {
			err1 = repo.Store(cipher)
			wg.Done()
		}()
		go func() {
			err2 = repo.Store(cipher)
			wg.Done()
		}()
		go func() {
			err3 = repo.Store(cipher)
			wg.Done()
		}()

		wg.Wait()

		Expect([]error{err1, err2, err3}).To(ConsistOf(
			BeNil(),
			MatchError("already exists"),
			MatchError("already exists"),
		))
	})

	It("does not read ciphers during writes", func() {
		repo := inmemory.New()

		cipher := repository.Cipher{
			ID:         "some-id",
			Nonce:      "some-nonce",
			CipherText: "some-cipher-text",
		}

		err := repo.Store(cipher)
		Expect(err).NotTo(HaveOccurred())

		wg := sync.WaitGroup{}
		wg.Add(3)

		var cipher1, cipher2 repository.Cipher
		var storeErr, err1, err2 error
		go func() {
			storeErr = repo.Store(repository.Cipher{ID: "another-id"})
			wg.Done()
		}()
		go func() {
			cipher1, err1 = repo.FindByID("some-id")
			wg.Done()
		}()
		go func() {
			cipher2, err2 = repo.FindByID("some-id")
			wg.Done()
		}()

		wg.Wait()

		Expect([]error{storeErr, err1, err2}).To(ConsistOf(
			BeNil(),
			BeNil(),
			BeNil(),
		))
		Expect([]repository.Cipher{cipher1, cipher2}).To(ConsistOf(
			cipher,
			cipher,
		))
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
