package repository

import "errors"

var (
	// ErrMarshaling is returned when a json result cannot be marshaled
	// into an entity.
	ErrMarshaling = errors.New("marshaling error")

	// ErrRequest is returned when a request to Elasticsearch server returns
	// a bad http status code.
	ErrRequest = errors.New("elasticsearch request error")

	// ErrResourceNotFound is returned when a query by ID has no match.
	ErrResourceNotFound = errors.New("resource not found")

	// ErrInternal is returned when an encountered error could not be identified.
	ErrInternal = errors.New("repository internal error")
)
