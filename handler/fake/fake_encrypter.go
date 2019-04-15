package fake

import "github.com/jamesjoshuahill/ciphers/encryption"

type FakeEncrypter struct {
	EncryptCall struct {
		Received struct {
			PlainText string
		}
		Returns struct {
			Cipher encryption.Cipher
			Error  error
		}
	}
}

func (e *FakeEncrypter) Encrypt(plainText string) (encryption.Cipher, error) {
	e.EncryptCall.Received.PlainText = plainText
	return e.EncryptCall.Returns.Cipher, e.EncryptCall.Returns.Error
}
