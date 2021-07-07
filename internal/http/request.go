package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// readURLQuery returns the value for the given key in the request URL.
// It returns a non nil error if the result is empty.
func (s Server) readURLQuery(r *http.Request, key string) (string, error) {
	q := r.URL.Query().Get(key)
	if q == "" {
		return "", errors.New("missing query: must use ?query=keywords in the url")
	}
	return q, nil
}

// extractRouteParam retreives the given route parameter from the
// mux path variables.
func extractRouteParam(r *http.Request, p string) (string, error) {
	v, ok := mux.Vars(r)[p]
	if !ok {
		return "", fmt.Errorf("invalid route parameter for \"%s\"", p)
	}
	return v, nil
}

// decodeBody reads the given request body and writes the decoded data to dest.
// The body is expected to be encoded as JSON.
func decodeBody(body io.ReadCloser, dest interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return err
	}
	return nil
}
