package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/handler"
)

type unexpectedError struct {
	statusCode int
	message    string
}

func newUnexpectedError(res *http.Response) unexpectedError {
	unerr := unexpectedError{statusCode: res.StatusCode}

	var body handler.ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return unerr
	}

	unerr.message = body.Message
	return unerr
}

func (e unexpectedError) Error() string {
	return fmt.Sprintf("status code: %d, message: %s", e.statusCode, e.message)
}

func (e unexpectedError) StatusCode() int {
	return e.statusCode
}

func (e unexpectedError) Message() string {
	return e.message
}

type alreadyExistsError struct{}

func (e alreadyExistsError) Error() string {
	return fmt.Sprintf("cipher already exists")
}

func (e alreadyExistsError) AlreadyExists() bool {
	return true
}

type notFoundError struct{}

func (e notFoundError) Error() string {
	return fmt.Sprintf("cipher not found")
}

func (e notFoundError) NotFound() bool {
	return true
}

type wrongKeyError struct{}

func (e wrongKeyError) Error() string {
	return fmt.Sprintf("wrong key for cipher")
}

func (e wrongKeyError) WrongKey() bool {
	return true
}
