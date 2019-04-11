package inmemory

import (
	"errors"
	"sync"

	"github.com/jamesjoshuahill/ciphers/repository"
)

type repo struct {
	ciphers map[string]repository.Cipher
	mutex   *sync.RWMutex
}

func New() *repo {
	return &repo{
		ciphers: make(map[string]repository.Cipher),
		mutex:   &sync.RWMutex{},
	}
}

func (r *repo) Store(cipher repository.Cipher) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.ciphers[cipher.ID]; ok {
		return errors.New("already exists")
	}

	r.ciphers[cipher.ID] = cipher

	return nil
}

func (r *repo) FindByID(id string) (repository.Cipher, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	cipher, ok := r.ciphers[id]
	if !ok {
		return repository.Cipher{}, errors.New("not found")
	}

	return cipher, nil
}
