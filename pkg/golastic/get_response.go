package golastic

import (
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// esSearchResponse represents the structure of an Elasticsearch response
// for a GET request.
type esGetResponse struct {
	Found  bool            `json:"found"`
	ID     string          `json:"_id"`
	Source json.RawMessage `json:"_source"`
}

// ReadGetResponse reads an Elasticsearch response for a GET request
// and returns the result as a generic interface or the first non-nil
// error occurring in the process.
//
// It must be provided a Document to determinate the marshaling process.
// The typical usage is to provide an entity having a custom NewHit method
// (see Document interface).
func ReadGetResponse(res *esapi.Response, doc Document) (interface{}, error) {
	if res.IsError() {
		return nil, statusError(res.StatusCode)
	}

	r, err := decodeRawGetResponse(res)
	if err != nil {
		return SearchResults{}, ErrUnhandled
	}

	return unmarshalGetResponse(r, doc)
}

func decodeRawGetResponse(res *esapi.Response) (esGetResponse, error) {
	var r esGetResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return r, err
	}
	return r, nil
}

func unmarshalGetResponse(r esGetResponse, doc Document) (interface{}, error) {
	if !r.Found {
		return nil, errors.New("not found")
	}

	result, err := doc.NewHit(r.ID, r.Source)
	if err != nil {
		return nil, err
	}

	return result, nil
}
