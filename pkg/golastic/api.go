package golastic

import (
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// ContextConfig configures the context for a Elasticsearch API call.
type ContextConfig struct {
	Client    *elasticsearch.Client
	IndexName string
}

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
// Deprecated: in the future this will be internally implemented
// in all ReadXxxResponse functions (such as ReadSearchResponse).
// Use for response kinds that do not have an associated ReadResponse function.
func ReadErrorResponse(res *esapi.Response) error {
	if !res.IsError() {
		return nil
	}
	return statusError(res.StatusCode)
}
