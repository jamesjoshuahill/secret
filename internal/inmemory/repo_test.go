package inmemory_test

import (
	"sync"

	"github.com/jamesjoshuahill/secret/internal/inmemory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	It("stores and retrieves secrets", func() {
		repo := inmemory.NewRepo()

		secret := inmemory.Secret{
			ID:         "some-id",
			Nonce:      "some-nonce",
			CipherText: "some-cipher-text",
		}
		err := repo.Store(secret)
		Expect(err).NotTo(HaveOccurred())

		actualSecret, err := repo.FindByID("some-id")
		Expect(err).NotTo(HaveOccurred())

		Expect(actualSecret).To(Equal(secret))
	})

	It("handles concurrent stores", func() {
		repo := inmemory.NewRepo()

		secret := inmemory.Secret{
			ID:         "some-id",
			Nonce:      "some-nonce",
			CipherText: "some-cipher-text",
		}

		wg := sync.WaitGroup{}
		wg.Add(3)

		var err1, err2, err3 error
		go func() {
			err1 = repo.Store(secret)
			wg.Done()
		}()
		go func() {
			err2 = repo.Store(secret)
			wg.Done()
		}()
		go func() {
			err3 = repo.Store(secret)
			wg.Done()
		}()

		wg.Wait()

		Expect([]error{err1, err2, err3}).To(ConsistOf(
			BeNil(),
			MatchError("already exists"),
			MatchError("already exists"),
		))
	})

	It("handles concurrent finds", func() {
		repo := inmemory.NewRepo()

		secret := inmemory.Secret{
			ID:         "some-id",
			Nonce:      "some-nonce",
			CipherText: "some-cipher-text",
		}

		err := repo.Store(secret)
		Expect(err).NotTo(HaveOccurred())

		wg := sync.WaitGroup{}
		wg.Add(3)

		var secret1, secret2 inmemory.Secret
		var storeErr, err1, err2 error
		go func() {
			storeErr = repo.Store(inmemory.Secret{ID: "another-id"})
			wg.Done()
		}()
		go func() {
			secret1, err1 = repo.FindByID("some-id")
			wg.Done()
		}()
		go func() {
			secret2, err2 = repo.FindByID("some-id")
			wg.Done()
		}()

		wg.Wait()

		Expect([]error{storeErr, err1, err2}).To(ConsistOf(
			BeNil(),
			BeNil(),
			BeNil(),
		))
		Expect([]inmemory.Secret{secret1, secret2}).To(ConsistOf(
			secret,
			secret,
		))
	})

	It("fails when the secret already exits", func() {
		repo := inmemory.NewRepo()

		err := repo.Store(inmemory.Secret{ID: "some-id"})
		Expect(err).NotTo(HaveOccurred())

		err = repo.Store(inmemory.Secret{ID: "some-id"})
		Expect(err).To(MatchError("already exists"))
	})

	It("fails when it cannot find a secret", func() {
		repo := inmemory.NewRepo()

		_, err := repo.FindByID("some-id")
		Expect(err).To(MatchError("not found"))
	})
})
