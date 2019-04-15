package handler

import (
	"github.com/jamesjoshuahill/ciphers/encryption"
	"github.com/jamesjoshuahill/ciphers/repository"
)

type Repository interface {
	Store(repository.Secret) error
	FindByID(string) (repository.Secret, error)
}

type Encrypter interface {
	Encrypt(string) (encryption.Secret, error)
}

type Decrypter interface {
	Decrypt(secret encryption.Secret) (string, error)
}
