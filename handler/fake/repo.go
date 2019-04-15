package fake

import (
	"github.com/jamesjoshuahill/ciphers/repository/inmemory"
)

type Repo struct {
	StoreCall struct {
		Received struct {
			Secret inmemory.Secret
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
			Secret inmemory.Secret
			Error  error
		}
	}
}

func (r *Repo) Store(secret inmemory.Secret) error {
	r.StoreCall.Received.Secret = secret
	return r.StoreCall.Returns.Error
}

func (r *Repo) FindByID(id string) (inmemory.Secret, error) {
	r.FindByResourceIDCall.Received.ID = id
	return r.FindByResourceIDCall.Returns.Secret, r.FindByResourceIDCall.Returns.Error
}
