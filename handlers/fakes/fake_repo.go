package fakes

import "github.com/jamesjoshuahill/ciphers/repository"

type FakeRepo struct {
	StoreCall struct {
		Received struct {
			Cipher repository.Cipher
		}
		Returns struct {
			Error error
		}
	}
}

func (r *FakeRepo) Store(cipher repository.Cipher) error {
	r.StoreCall.Received.Cipher = cipher
	return r.StoreCall.Returns.Error
}
