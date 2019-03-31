package fakes

type FakeEncrypter struct {
	EncryptCall struct {
		Received struct {
			Plaintext string
		}
		Returns struct {
			Key        string
			Ciphertext string
			Error      error
		}
	}
}

func (e *FakeEncrypter) Encrypt(plaintext string) (string, string, error) {
	e.EncryptCall.Received.Plaintext = plaintext
	return e.EncryptCall.Returns.Key, e.EncryptCall.Returns.Ciphertext, e.EncryptCall.Returns.Error
}
