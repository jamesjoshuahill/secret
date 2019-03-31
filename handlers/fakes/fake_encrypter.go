package fakes

type FakeEncrypter struct {
	EncryptCall struct {
		Received struct {
			PlainText string
		}
		Returns struct {
			Key        string
			CipherText string
			Error      error
		}
	}
}

func (e *FakeEncrypter) Encrypt(plainText string) (string, string, error) {
	e.EncryptCall.Received.PlainText = plainText
	return e.EncryptCall.Returns.Key, e.EncryptCall.Returns.CipherText, e.EncryptCall.Returns.Error
}
