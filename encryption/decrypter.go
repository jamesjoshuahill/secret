package encryption

type Decrypter struct{}

func (Decrypter) Decrypt(key, cipherText string) (string, error) {
	return "some plain text", nil
}
