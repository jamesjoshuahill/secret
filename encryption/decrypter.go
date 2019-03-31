package encryption

import "errors"

type Decrypter struct{}

func (Decrypter) Decrypt(key, cipherText string) (string, error) {
	if key != hexKey {
		return "", errors.New("wrong key")
	}

	return "some plain text", nil
}
