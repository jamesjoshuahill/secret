package handler

import (
	"github.com/jamesjoshuahill/ciphers/aes"
	"github.com/jamesjoshuahill/ciphers/inmemory"
)

type Repository interface {
	Store(inmemory.Secret) error
	FindByID(string) (inmemory.Secret, error)
}

type Encrypter interface {
	Encrypt(string) (aes.Secret, error)
}

type Decrypter interface {
	Decrypt(secret aes.Secret) (string, error)
}
