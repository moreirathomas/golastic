package golastic

import (
	"errors"
	"net/http"
)

var (
	// ErrBadRequest is returned when a request is malformed
	ErrBadRequest = errors.New("bad request")

	// ErrNotFound is returned when a requested resource is not found.
	ErrNotFound = errors.New("resource not found")

	// ErrUnhandled is returned when an encountered error cannot be identified.
	ErrUnhandled = errors.New("elasticsearch unhandled error")
)

var statusErrorMapping = map[int]error{
	http.StatusBadRequest:          ErrBadRequest,
	http.StatusNotFound:            ErrNotFound,
	http.StatusInternalServerError: ErrUnhandled,
}

func statusError(code int) error {
	if err := statusErrorMapping[code]; err != nil {
		return err
	}
	return ErrUnhandled
}
