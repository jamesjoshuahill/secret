package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

type Encrypter struct{}

func (Encrypter) Encrypt(plaintext string) (Cipher, error) {
	secretKey := make([]byte, 32)
	_, err := rand.Read(secretKey)
	if err != nil {
		return Cipher{}, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return Cipher{}, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return Cipher{}, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return Cipher{}, err
	}

	cipherText := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	return Cipher{
		Key:        hex.EncodeToString(secretKey),
		Nonce:      hex.EncodeToString(nonce),
		CipherText: hex.EncodeToString(cipherText),
	}, nil
}
