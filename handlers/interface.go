package handlers

import (
	"github.com/jamesjoshuahill/ciphers/encryption"
	"github.com/jamesjoshuahill/ciphers/repository"
)

type Repository interface {
	Store(repository.Cipher) error
	FindByID(string) (repository.Cipher, error)
}

type Encrypter interface {
	Encrypt(string) (encryption.Cipher, error)
}

type Decrypter interface {
	Decrypt(cipher encryption.Cipher) (string, error)
}
