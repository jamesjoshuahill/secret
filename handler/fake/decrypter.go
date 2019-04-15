package fake

import "github.com/jamesjoshuahill/ciphers/aes"

type Decrypter struct {
	DecryptCall struct {
		Received struct {
			Secret aes.Secret
		}
		Returns struct {
			PlainText string
			Error     error
		}
	}
}

func (d *Decrypter) Decrypt(secret aes.Secret) (string, error) {
	d.DecryptCall.Received.Secret = secret
	return d.DecryptCall.Returns.PlainText, d.DecryptCall.Returns.Error
}
