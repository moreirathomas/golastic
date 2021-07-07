package http

import (
	"fmt"
	"net/http"
)

var (
	ErrBadRequest   = HTTPError{nil, http.StatusText(http.StatusBadRequest), http.StatusBadRequest}
	ErrUnauthorized = HTTPError{nil, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized}
	ErrNotFound     = HTTPError{nil, http.StatusText(http.StatusNotFound), http.StatusNotFound}
	ErrInternal     = HTTPError{nil, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError}
)

// HTTPError is a high-level error that wraps another error
// and holds a HTTP status code.
type HTTPError struct {
	wrapped error  // wrapped is used internally to share the error context.
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error returns an error's message.
func (e HTTPError) Error() string {
	return e.Message
}

// Wrap wraps the given error, keeping the receiver's status code.
func (e HTTPError) Wrap(target error) HTTPError {
	return HTTPError{
		Code:    e.Code,
		Message: fmt.Sprintf("%s: %s", e.Message, target.Error()),
		wrapped: target,
	}
}

// Unwrap returns the receiver's wrapped error.
func (e HTTPError) Unwrap() error {
	return e.wrapped
}
