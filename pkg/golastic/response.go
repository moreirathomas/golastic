package golastic

import (
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// Hit represents a single result as returned by an Elasticsearch response.
type Hit struct {
	ID     string          `json:"_id"`
	Score  float32         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

// Unmarshaler expects an UnmarshalHit method that is used to unmarshal a Hit
// retrieved from the body of an Elasticsearch API response.
type Unmarshaler interface {
	// UnmarshalHit returns the document as an unmarshalled struct
	// or a non-nil error.
	UnmarshalHit(Hit) (interface{}, error)
}

// ReadErrorResponse reads the response body and returns an error if
// the response status indicates failure.
func ReadErrorResponse(res *esapi.Response) error {
	if !res.IsError() {
		return nil
	}
	return statusError(res.StatusCode)
}
