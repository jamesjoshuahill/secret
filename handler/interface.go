package handler

import (
	"github.com/jamesjoshuahill/secret/aes"
	"github.com/jamesjoshuahill/secret/inmemory"
)

type Repository interface {
	Store(inmemory.Secret) error
	FindByID(string) (inmemory.Secret, error)
}

type EncryptFunc func(string) (aes.Secret, error)
type DecryptFunc func(secret aes.Secret) (string, error)
