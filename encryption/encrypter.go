package encryption

type Encrypter struct{}

func (Encrypter) Encrypt(plaintext string) (string, string, error) {
	return "key for client-cipher-id", "", nil
}
