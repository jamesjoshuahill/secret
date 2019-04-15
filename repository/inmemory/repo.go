package inmemory

import (
	"errors"
	"sync"

	"github.com/jamesjoshuahill/ciphers/repository"
)

type repo struct {
	mutex   *sync.RWMutex
	secrets map[string]repository.Secret
}

func NewRepo() *repo {
	return &repo{
		mutex:   &sync.RWMutex{},
		secrets: make(map[string]repository.Secret),
	}
}

func (r *repo) Store(secret repository.Secret) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.secrets[secret.ID]; ok {
		return errors.New("already exists")
	}

	r.secrets[secret.ID] = secret

	return nil
}

func (r *repo) FindByID(id string) (repository.Secret, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	secret, ok := r.secrets[id]
	if !ok {
		return repository.Secret{}, errors.New("not found")
	}

	return secret, nil
}
