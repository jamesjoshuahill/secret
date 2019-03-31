package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const ivHex = "64a9433eae7ccceee2fc0eda"

type Encrypter struct{}

func (Encrypter) Encrypt(plaintext string) (string, string, error) {
	secretKey := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, secretKey)
	if err != nil {
		return "", "", err
	}

	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	cipherText := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return hex.EncodeToString(secretKey), hex.EncodeToString(cipherText), nil
}
