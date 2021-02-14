package inmemory

import (
	"errors"
	"sync"
)

type Repo struct {
	mutex   *sync.RWMutex
	secrets map[string]Secret
}

func NewRepo() *Repo {
	return &Repo{
		mutex:   &sync.RWMutex{},
		secrets: make(map[string]Secret),
	}
}

func (r *Repo) Store(secret Secret) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.secrets[secret.ID]; ok {
		return errors.New("already exists")
	}

	r.secrets[secret.ID] = secret

	return nil
}

func (r *Repo) FindByID(id string) (Secret, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	secret, ok := r.secrets[id]
	if !ok {
		return Secret{}, errors.New("not found")
	}

	return secret, nil
}
