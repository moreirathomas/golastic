package golastic

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// ContextConfig configures the context for a Elasticsearch API call.
type ContextConfig struct {
	Client    *elasticsearch.Client
	IndexName string
}

// Document provides a NewHit method that is used to unmarshal a json
// from the body of an Elasticsearch API response.
type Document interface {
	// NewHit returns the document as an unmarshalled struct
	// or a non-nil error.
	NewHit(id string, src json.RawMessage) (interface{}, error)
}

// ReadErrorResponse reads the response body and returns an error if
// the response status indicates failure.
func ReadErrorResponse(res *esapi.Response) error {
	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res)
	}
	return nil
}
