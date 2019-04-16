package handler

import (
	"github.com/jamesjoshuahill/ciphers/aes"
	"github.com/jamesjoshuahill/ciphers/inmemory"
)

type Repository interface {
	Store(inmemory.Secret) error
	FindByID(string) (inmemory.Secret, error)
}

type EncryptFunc func(string) (aes.Secret, error)
type DecryptFunc func(secret aes.Secret) (string, error)
