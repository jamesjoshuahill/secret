package client

type client struct{}

func New() *client {
	return &client{}
}

func (*client) Store(id, payload []byte) ([]byte, error) {
	return nil, nil
}

func (*client) Retrieve(id, aesKey []byte) ([]byte, error) {
	return nil, nil
}
