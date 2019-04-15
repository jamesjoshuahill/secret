package encryption_test

import (
	"github.com/jamesjoshuahill/ciphers/encryption"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decrypter", func() {
	It("decrypts cipher text", func() {
		decrypter := encryption.Decrypter{}

		plainText, err := decrypter.Decrypt(encryption.Secret{
			Key:        "6368616e676520746869732070617373776f726420746f206120736563726574",
			Nonce:      "64a9433eae7ccceee2fc0eda",
			CipherText: "c3aaa29f002ca75870806e44086700f62ce4d43e902b3888e23ceff797a7a471",
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(plainText).To(Equal("exampleplaintext"))
	})
})
