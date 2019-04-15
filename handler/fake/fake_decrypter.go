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

func (d *FakeDecrypter) Decrypt(secret encryption.Secret) (string, error) {
	d.DecryptCall.Received.Secret = secret
	return d.DecryptCall.Returns.PlainText, d.DecryptCall.Returns.Error
}
