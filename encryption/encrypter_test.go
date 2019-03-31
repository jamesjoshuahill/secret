package encryption_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jamesjoshuahill/ciphers/encryption"
)

var _ = Describe("Encrypter", func() {
	It("encrypts plain text", func() {
		encrypter := encryption.Encrypter{}

		_, cipherText1, err := encrypter.Encrypt("some plain text")
		Expect(err).NotTo(HaveOccurred())

		_, cipherText2, err := encrypter.Encrypt("some plain text")
		Expect(err).NotTo(HaveOccurred())

		Expect(cipherText2).NotTo(Equal(cipherText1))
	})
})
