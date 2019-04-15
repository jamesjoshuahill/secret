package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

type Decrypter struct{}

func (Decrypter) Decrypt(s Secret) (string, error) {
	secretKey, err := hex.DecodeString(s.Key)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(s.CipherText)
	if err != nil {
		return "", err
	}

	nonce, err := hex.DecodeString(s.Nonce)
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
