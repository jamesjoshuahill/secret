package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/handler"
)

var ErrWrongIDOrKey = errors.New("wrong id or key")

type UnexpectedError struct {
	StatusCode int
	Message    string
}

func newUnexpectedError(res *http.Response) UnexpectedError {
	unerr := UnexpectedError{StatusCode: res.StatusCode}

	var body handler.ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return unerr
	}

	unerr.Message = body.Message
	return unerr
}

func (e UnexpectedError) Error() string {
	return fmt.Sprintf("status code: %d, message: %s", e.StatusCode, e.Message)
}

type alreadyExistsError struct{}

func (e alreadyExistsError) Error() string {
	return fmt.Sprintf("secret already exists")
}

func (e alreadyExistsError) AlreadyExists() bool {
	return true
}
