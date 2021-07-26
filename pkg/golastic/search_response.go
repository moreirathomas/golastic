package golastic

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// SearchResults is a simplified and flattened result
// of an Elasticsearch response for a search query.
type SearchResults struct {
	Total int           `json:"total"`
	Hits  []interface{} `json:"hits"`
}

// esSearchResponse represents the structure of an Elasticsearch response
// for a search query.
type esSearchResponse struct {
	Took int
	Hits struct {
		Total struct {
			Value int
		}
		Hits []struct {
			ID     string          `json:"_id"`
			Source json.RawMessage `json:"_source"`
			Sort   []interface{}   `json:"sort"`
		}
	}
}

// ReadSearchResponse reads an Elasticsearch response for a Search request
// and returns a SearchResults or the first non-nil error occurring
// in the process.
//
// It must be provided a Document to determinate the marshaling process.
// The typical usage is to provide an entity having a custom NewHit method
// (see Document interface).
func ReadSearchResponse(res *esapi.Response, doc Document) (SearchResults, error) {
	if res.IsError() {
		return SearchResults{}, statusError(res.StatusCode)
	}

	r, err := decodeRawSearchResponse(res)
	if err != nil {
		return SearchResults{}, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	return unmarshalSearchResponse(r, doc)
}

func decodeRawSearchResponse(res *esapi.Response) (esSearchResponse, error) {
	var r esSearchResponse
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return r, err
	}
	return r, nil
}

func unmarshalSearchResponse(r esSearchResponse, doc Document) (SearchResults, error) {
	var results SearchResults

	results.Total = r.Hits.Total.Value

	// handle empty results
	if len(r.Hits.Hits) < 1 {
		results.Hits = []interface{}{}
		return results, nil
	}

	// unmarshal elasticsearch hits
	for _, hit := range r.Hits.Hits {
		h, err := doc.NewHit(hit.ID, hit.Source)
		if err != nil {
			return results, err
		}
		results.Hits = append(results.Hits, h)
	}

	return results, nil
}
