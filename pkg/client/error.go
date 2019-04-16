package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/secret/handler"
)

var (
	ErrWrongIDOrKey  = errors.New("wrong id or key")
	ErrAlreadyExists = errors.New("secret already exists")
)

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
	err := json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return unerr
	}

	unerr.Message = body.Message
	return unerr
}
