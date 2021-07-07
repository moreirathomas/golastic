package http

import (
	"encoding/json"
	"net/http"
)

// setHeader is a helper function for writing response header' content type and code.
func setHeader(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}

// respondJSON sends the given data as JSON. The response status code is set to the given code.
func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	setHeader(w, code)

	resp, err := json.Marshal(data)
	if err != nil {
		respondHTTPError(w, ErrInternal)
		return
	}
	w.Write(resp)
}

// respondHTTPError formats the given error and sends it as JSON.
// Parameter `httpErr` sets the status code.
func respondHTTPError(w http.ResponseWriter, httpErr HTTPError) {
	resp := map[string]HTTPError{
		"error": httpErr,
	}
	respondJSON(w, httpErr.Code, resp)
}
