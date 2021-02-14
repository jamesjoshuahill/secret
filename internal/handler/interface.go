package handler

import (
	"github.com/jamesjoshuahill/secret/internal/aes"
	"github.com/jamesjoshuahill/secret/internal/inmemory"
)

type Repository interface {
	Store(inmemory.Secret) error
	FindByID(string) (inmemory.Secret, error)
}

type (
	EncryptFunc func(string) (aes.Secret, error)
	DecryptFunc func(secret aes.Secret) (string, error)
)
