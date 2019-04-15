package fake

import "net/http"

type HTTPSClient struct {
	DoCall struct {
		Received struct {
			Request *http.Request
		}
		Returns struct {
			Response *http.Response
			Error    error
		}
	}
}

func (f *HTTPSClient) Do(req *http.Request) (*http.Response, error) {
	f.DoCall.Received.Request = req
	return f.DoCall.Returns.Response, f.DoCall.Returns.Error
}
