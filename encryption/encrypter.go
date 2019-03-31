package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const hexKey = "6368616e676520746869732070617373776f726420746f206120736563726574"

type Encrypter struct{}

func (Encrypter) Encrypt(plaintext string) (string, string, error) {
	key, _ := hex.DecodeString(hexKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	cipherText := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(nonce), hex.EncodeToString(cipherText), nil
}
