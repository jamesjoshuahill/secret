package repository

type repo struct{}

func New() *repo {
	return &repo{}
}

func (r *repo) Store(cipher Cipher) error {
	return nil
}

func (r *repo) FindByResourceID(resourceID string) (Cipher, error) {
	return Cipher{
		ResourceID: resourceID,
		CipherText: "some plain text",
		Key:        "key for client-cipher-id",
	}, nil
}
