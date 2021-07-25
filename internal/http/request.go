package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/moreirathomas/golastic/internal"
)

// extractQueryParam returns the given param value in the request query.
func extractQueryParam(r *http.Request, p string) string {
	return r.URL.Query().Get(p)
}

// extractQueryParamInt returns the given param value in the request query.
// It returns a non nil error if the result is not a number.
func extractQueryParamInt(r *http.Request, p string) (int, error) {
	qStr := extractQueryParam(r, p)
	q, err := strconv.Atoi(qStr)
	if err != nil {
		return 0, fmt.Errorf("bad or missing query parameter: \"%s\" must be a number", p)
	}
	return q, nil
}

// extractRouteParam retreives the given route parameter from the
// mux path variables.
func extractRouteParam(r *http.Request, p string) (string, error) {
	v, ok := mux.Vars(r)[p]
	if !ok {
		return "", fmt.Errorf("bad route parameter for \"%s\"", p)
	}
	return v, nil
}

// decodeBody reads the given request body and writes the decoded data to dest.
// The body is expected to be encoded as JSON.
func decodeBody(body io.ReadCloser, dest interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dest)
}

func readBookPayload(body io.ReadCloser) (internal.Book, error) {
	var book internal.Book

	if err := decodeBody(body, &book); err != nil {
		return internal.Book{}, err
	}

	if err := book.Validate(); err != nil {
		return internal.Book{}, err
	}

	return book, nil
}
