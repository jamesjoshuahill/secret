package fakes

import "github.com/jamesjoshuahill/ciphers/encryption"

type FakeDecrypter struct {
	DecryptCall struct {
		Received struct {
			Cipher encryption.Cipher
		}
		Returns struct {
			PlainText string
			Error     error
		}
	}
}

func (d *FakeDecrypter) Decrypt(cipher encryption.Cipher) (string, error) {
	d.DecryptCall.Received.Cipher = cipher
	return d.DecryptCall.Returns.PlainText, d.DecryptCall.Returns.Error
}
