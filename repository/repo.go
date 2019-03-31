package repository

type repo struct{}

func New() *repo {
	return &repo{}
}

func (r repo) Store(cipher Cipher) error {
	return nil
}
