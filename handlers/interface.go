package handlers

import "github.com/jamesjoshuahill/ciphers/repository"

type Repository interface {
	Store(repository.Cipher) error
}
