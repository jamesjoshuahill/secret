package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

type Decrypter struct{}

func (Decrypter) Decrypt(c Cipher) (string, error) {
	secretKey, err := hex.DecodeString(c.Key)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(c.CipherText)
	if err != nil {
		return "", err
	}

	nonce, err := hex.DecodeString(c.Nonce)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainText, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
