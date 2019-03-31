package repository

type repo struct{}

func New() *repo {
	return &repo{}
}

func (r *repo) Store(cipher Cipher) error {
	return nil
}

func (r *repo) FindByID(id string) (Cipher, error) {
	return Cipher{
		ID:         id,
		CipherText: "some cipher text",
	}, nil
}
