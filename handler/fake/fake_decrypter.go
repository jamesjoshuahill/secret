package fake

import "github.com/jamesjoshuahill/ciphers/encryption"

type FakeDecrypter struct {
	DecryptCall struct {
		Received struct {
			Secret encryption.Secret
		}
		Returns struct {
			PlainText string
			Error     error
		}
	}
}

func (d *FakeDecrypter) Decrypt(cipher encryption.Secret) (string, error) {
	d.DecryptCall.Received.Secret = cipher
	return d.DecryptCall.Returns.PlainText, d.DecryptCall.Returns.Error
}
