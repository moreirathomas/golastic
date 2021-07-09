package golastic

import (
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// Document provides a NewHit method that should unmarshal a json result
// obtained from an ElasticSearch search response.
// It should return your own entity and an error.
type Document interface {
	NewHit(id string, src json.RawMessage) (interface{}, error)
}

// SearchResults is a simplified result of an ElasticSearch output.
type SearchResults struct {
	Total int           `json:"total"`
	Hits  []interface{} `json:"hits"`
}

// searchResponseWrapper represents the response of an ElasticSearch search query.
type searchResponseWrapper struct {
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

// UnwrapSearchResponse reads an ElasticSearch response and returns a SearchResults
// or the first non-nil error occurring in the process.
// It must be provided a Document to determinate the marshaling process.
// The typical usage is to provide an entity having a custom NewHit method
// (see Document interface).
func UnwrapSearchResponse(res *esapi.Response, doc Document) (SearchResults, error) {
	var results SearchResults

	var rw searchResponseWrapper
	if err := json.NewDecoder(res.Body).Decode(&rw); err != nil {
		return results, err
	}

	results.Total = rw.Hits.Total.Value
	if len(rw.Hits.Hits) < 1 {
		results.Hits = []interface{}{}
		return results, nil
	}

	for _, hit := range rw.Hits.Hits {
		h, err := doc.NewHit(hit.ID, hit.Source)
		if err != nil {
			return SearchResults{}, err
		}
		results.Hits = append(results.Hits, h)
	}

	return results, nil
}
