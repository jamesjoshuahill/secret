package fake

import "github.com/jamesjoshuahill/ciphers/aes"

type Encrypter struct {
	EncryptCall struct {
		Received struct {
			PlainText string
		}
		Returns struct {
			Secret aes.Secret
			Error  error
		}
	}
}

func (e *Encrypter) Encrypt(plainText string) (aes.Secret, error) {
	e.EncryptCall.Received.PlainText = plainText
	return e.EncryptCall.Returns.Secret, e.EncryptCall.Returns.Error
}
