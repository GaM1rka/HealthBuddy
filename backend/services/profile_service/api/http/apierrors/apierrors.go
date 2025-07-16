package apierrors

import "fmt"

type APIError struct {
	Code    int
	Message string
	Err     error
}

func (e APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%d: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func NewBadRequest(msg string, err error) APIError {
	return APIError{Code: 400, Message: msg, Err: err}
}
func NewNotFound(msg string) APIError {
	return APIError{Code: 404, Message: msg}
}
func NewForbidden(msg string) APIError {
	return APIError{Code: 403, Message: msg}
}
func NewInternal(err error) APIError {
	return APIError{Code: 500, Message: "internal error", Err: err}
}
