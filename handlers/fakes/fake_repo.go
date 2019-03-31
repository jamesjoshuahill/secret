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
	FindByResourceIDCall struct {
		Received struct {
			ID string
		}
		Returns struct {
			Cipher repository.Cipher
			Error  error
		}
	}
}

func (r *FakeRepo) Store(cipher repository.Cipher) error {
	r.StoreCall.Received.Cipher = cipher
	return r.StoreCall.Returns.Error
}

func (r *FakeRepo) FindByID(id string) (repository.Cipher, error) {
	r.FindByResourceIDCall.Received.ID = id
	return r.FindByResourceIDCall.Returns.Cipher, r.FindByResourceIDCall.Returns.Error
}
