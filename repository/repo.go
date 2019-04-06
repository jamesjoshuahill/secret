package repository

import "errors"

type repo struct {
	ciphers map[string]Cipher
}

func New() *repo {
	return &repo{
		ciphers: make(map[string]Cipher),
	}
}

func (r *repo) Store(cipher Cipher) error {
	if _, ok := r.ciphers[cipher.ID]; ok {
		return errors.New("already exists")
	}

	r.ciphers[cipher.ID] = cipher

	return nil
}

func (r *repo) FindByID(id string) (Cipher, error) {
	cipher, ok := r.ciphers[id]
	if !ok {
		return Cipher{}, errors.New("not found")
	}

	return cipher, nil
}
