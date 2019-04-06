package encryption_test

import (
	"encoding/hex"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jamesjoshuahill/ciphers/encryption"
)

var _ = Describe("Encrypter", func() {
	It("encrypts plain text", func() {
		encrypter := encryption.Encrypter{}

		cipher, err := encrypter.Encrypt("some plain text")
		Expect(err).NotTo(HaveOccurred())
		Expect(hex.DecodeString(cipher.Key)).To(HaveLen(32))
		Expect(hex.DecodeString(cipher.Nonce)).To(HaveLen(12))
	})

	It("creates a new key and nonce for each cipher", func() {
		encrypter := encryption.Encrypter{}

		cipher1, err := encrypter.Encrypt("some plain text")
		Expect(err).NotTo(HaveOccurred())
		cipher2, err := encrypter.Encrypt("some plain text")
		Expect(err).NotTo(HaveOccurred())

		Expect(cipher1.Key).NotTo(Equal(cipher2.Key))
		Expect(cipher1.Nonce).NotTo(Equal(cipher2.Nonce))
		Expect(cipher1.CipherText).NotTo(Equal(cipher2.CipherText))
	})
})
