package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

type Decrypter struct{}

func (Decrypter) Decrypt(key, cipherText string) (string, error) {
	secretKey, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	iv, err := hex.DecodeString(ivHex)
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

	plainText, err := aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
