package handlers

import "github.com/jamesjoshuahill/ciphers/repository"

type Repository interface {
	Store(repository.Cipher) error
	FindByResourceID(string) (repository.Cipher, error)
}

type Encrypter interface {
	Encrypt(string) (string, string, error)
}

type Decrypter interface {
	Decrypt(string, string) (string, error)
}
