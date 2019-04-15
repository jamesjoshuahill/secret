package aes_test

import (
	"encoding/hex"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jamesjoshuahill/ciphers/aes"
)

var _ = Describe("Encrypter", func() {
	It("encrypts plain text", func() {
		encrypter := aes.Encrypter{}

		secret, err := encrypter.Encrypt("some plain text")
		Expect(err).NotTo(HaveOccurred())
		Expect(hex.DecodeString(secret.Key)).To(HaveLen(32))
		Expect(hex.DecodeString(secret.Nonce)).To(HaveLen(12))
	})

	It("creates a new key and nonce for each secret", func() {
		encrypter := aes.Encrypter{}

		secret1, err := encrypter.Encrypt("some plain text")
		Expect(err).NotTo(HaveOccurred())
		secret2, err := encrypter.Encrypt("some plain text")
		Expect(err).NotTo(HaveOccurred())

		Expect(secret1.Key).NotTo(Equal(secret2.Key))
		Expect(secret1.Nonce).NotTo(Equal(secret2.Nonce))
		Expect(secret1.CipherText).NotTo(Equal(secret2.CipherText))
	})
})
