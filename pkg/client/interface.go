package client

import "net/http"

type HTTPSClient interface {
	Do(req *http.Request) (*http.Response, error)
}
