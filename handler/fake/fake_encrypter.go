package fake

import "github.com/jamesjoshuahill/ciphers/encryption"

type FakeEncrypter struct {
	EncryptCall struct {
		Received struct {
			PlainText string
		}
		Returns struct {
			Secret encryption.Secret
			Error  error
		}
	}
}

func (e *FakeEncrypter) Encrypt(plainText string) (encryption.Secret, error) {
	e.EncryptCall.Received.PlainText = plainText
	return e.EncryptCall.Returns.Secret, e.EncryptCall.Returns.Error
}
