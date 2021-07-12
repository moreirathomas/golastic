package http

import (
	"fmt"
	"net/http"
)

var (
	errBadRequest = httpError{nil, http.StatusText(http.StatusBadRequest), http.StatusBadRequest}
	errNotFound   = httpError{nil, http.StatusText(http.StatusNotFound), http.StatusNotFound}
	errInternal   = httpError{nil, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError}
)

// httpError is a high-level error that wraps another error
// and holds a HTTP status code.
type httpError struct {
	wrapped error  // wrapped is used internally to share the error context.
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error returns an error's message.
func (e httpError) Error() string {
	return e.Message
}

// Wrap wraps the given error, keeping the receiver's status code.
func (e httpError) Wrap(target error) httpError {
	return httpError{
		Code:    e.Code,
		Message: fmt.Sprintf("%s: %s", e.Message, target.Error()),
		wrapped: target,
	}
}

// Unwrap returns the receiver's wrapped error.
func (e httpError) Unwrap() error {
	return e.wrapped
}
