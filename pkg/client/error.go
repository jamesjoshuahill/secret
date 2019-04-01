package client

import "fmt"

type unexpectedError struct {
	statusCode int
	message    string
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
