package fakes

type FakeDecrypter struct {
	DecryptCall struct {
		Received struct {
			Key        string
			CipherText string
		}
		Returns struct {
			PlainText string
			Error     error
		}
	}
}

func (d *FakeDecrypter) Decrypt(key, cipherText string) (string, error) {
	d.DecryptCall.Received.Key = key
	d.DecryptCall.Received.CipherText = cipherText
	return d.DecryptCall.Returns.PlainText, d.DecryptCall.Returns.Error
}
