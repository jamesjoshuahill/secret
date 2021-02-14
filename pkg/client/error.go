package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/secret/internal/handler"
)

var (
	// ErrWrongIDOrKey is a sentinal error returned by Client.Retrieve().
	ErrWrongIDOrKey = errors.New("wrong id or key")
	// ErrAlreadyExists is a sentinal error returned by Client.Retrieve().
	ErrAlreadyExists = errors.New("secret already exists")
)

// UnexpectedError is an unexpected response from the server.
type UnexpectedError struct {
	StatusCode int
	Message    string
}

func (e UnexpectedError) Error() string {
	return fmt.Sprintf("status code: %d, message: %s", e.StatusCode, e.Message)
}

func newUnexpectedError(res *http.Response) UnexpectedError {
	unerr := UnexpectedError{StatusCode: res.StatusCode}

	var body handler.ErrorResponse
	defer res.Body.Close() //nolint:errcheck
	err := json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return unerr
	}

	unerr.Message = body.Message
	return unerr
}
