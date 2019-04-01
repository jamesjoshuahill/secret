package client

import "net/http"

// HTTPSClient is an interface representing the ability to make an http.Request
// and return the http.Response.
type HTTPSClient interface {
	Do(req *http.Request) (*http.Response, error)
}
