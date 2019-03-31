package repository

import "errors"

type repo struct {
	ciphers map[string]string
}

func New() *repo {
	return &repo{
		ciphers: make(map[string]string),
	}
}

func (r *repo) Store(cipher Cipher) error {
	r.ciphers[cipher.ID] = cipher.CipherText

	return nil
}

func (r *repo) FindByID(id string) (Cipher, error) {
	cipherText, ok := r.ciphers[id]
	if !ok {
		return Cipher{}, errors.New("not found")
	}

	return Cipher{
		ID:         id,
		CipherText: cipherText,
	}, nil
}
