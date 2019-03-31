package encryption_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jamesjoshuahill/ciphers/encryption"
)

var _ = Describe("Encrypter", func() {
	It("encrypts plain text", func() {
		encrypter := encryption.Encrypter{}

		key1, cipherText1, err := encrypter.Encrypt("some plain text")
		Expect(err).NotTo(HaveOccurred())

		key2, cipherText2, err := encrypter.Encrypt("some plain text")
		Expect(err).NotTo(HaveOccurred())

		Expect(key1).NotTo(Equal(key2))
		Expect(cipherText2).NotTo(Equal(cipherText1))
	})
})
