package handler

import (
	"github.com/jamesjoshuahill/ciphers/aes"
	"github.com/jamesjoshuahill/ciphers/repository"
)

type Repository interface {
	Store(repository.Secret) error
	FindByID(string) (repository.Secret, error)
}

type Encrypter interface {
	Encrypt(string) (aes.Secret, error)
}

type Decrypter interface {
	Decrypt(secret aes.Secret) (string, error)
}
