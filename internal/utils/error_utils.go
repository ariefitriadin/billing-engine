package utils

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func NewHTTPError(code int, message string) error {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

func ErrorCode(err error) int {
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr.Code
	}

	return http.StatusInternalServerError
}
