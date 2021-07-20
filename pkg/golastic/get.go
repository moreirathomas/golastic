package golastic

import (
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// getResponseWrapper represents selected fields from
// the response to an Elasticsearch GET request.
type getResponseWrapper struct {
	Found  bool            `json:"found"`
	ID     string          `json:"_id"`
	Source json.RawMessage `json:"_source"`
}

// UnwrapGetResponse reads an Elasticsearch response for a GET request
// and returns a generic interface or the first non-nil error occurring
// in the process.
//
// It must be provided a Document to determinate the marshaling process.
// The typical usage is to provide an entity having a custom NewHit method
// (see Document interface).
func UnwrapGetResponse(res *esapi.Response, doc Document) (interface{}, error) {
	var rw getResponseWrapper
	if err := json.NewDecoder(res.Body).Decode(&rw); err != nil {
		return rw, err
	}

	result, err := doc.NewHit(rw.ID, rw.Source)
	if err != nil {
		return nil, err
	}

	return result, nil
}
