package golastic

import (
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// esSearchResponse represents the structure of an Elasticsearch response
// for a GET request.
type esGetResponse struct {
	Found bool `json:"found"`
	Hit
}

// unwrap returns the document wrapped in esGetResponse
// or the first non-nil error encountered in the process.
// It uses Unmarshaler interface to determinate the marshaling process.
func (r esGetResponse) unwrap(doc Unmarshaler) (interface{}, error) {
	if !r.Found {
		return nil, errors.New("not found")
	}

	result, err := doc.UnmarshalHit(r.Hit)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ReadGetResponse reads an Elasticsearch response for a GET request
// and returns the result as a generic interface or the first non-nil
// error occurring in the process.
//
// It must be provided a Document to determinate the marshaling process.
// The typical usage is to provide an entity having a custom NewHit method
// (see Document interface).
func ReadGetResponse(res *esapi.Response, doc Unmarshaler) (interface{}, error) {
	if res.IsError() {
		return nil, statusError(res.StatusCode)
	}

	r, err := decodeRawGetResponse(res)
	if err != nil {
		return SearchResults{}, ErrUnhandled
	}

	return r.unwrap(doc)
}

func decodeRawGetResponse(res *esapi.Response) (esGetResponse, error) {
	var r esGetResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return r, err
	}
	return r, nil
}
