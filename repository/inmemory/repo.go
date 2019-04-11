package inmemory

import (
	"errors"

	"github.com/jamesjoshuahill/ciphers/repository"
)

type repo struct {
	ciphers map[string]repository.Cipher
}

func New() *repo {
	return &repo{
		ciphers: make(map[string]repository.Cipher),
	}
}

func (r *repo) Store(cipher repository.Cipher) error {
	if _, ok := r.ciphers[cipher.ID]; ok {
		return errors.New("already exists")
	}

	r.ciphers[cipher.ID] = cipher

	return nil
}

func (r *repo) FindByID(id string) (repository.Cipher, error) {
	cipher, ok := r.ciphers[id]
	if !ok {
		return repository.Cipher{}, errors.New("not found")
	}

	return cipher, nil
}
