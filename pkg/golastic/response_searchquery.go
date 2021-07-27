package golastic

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// SearchResults is a simplified and flattened result
// of an Elasticsearch response for a search query.
type SearchResults struct {
	Total   int
	Results []interface{}
}

// esSearchResponse represents the structure of an Elasticsearch response
// for a search query.
type esSearchResponse struct {
	Took int
	Hits struct {
		Total struct {
			Value int
		}
		Hits []Hit
	}
}

// unwrap returns the unwrap of esSearchResponse as SearchResults
// or the first non-nil error encountered in the process.
// It uses Unmarshaler interface to determinate the marshaling process.
func (r esSearchResponse) unwrap(doc Unmarshaler) (SearchResults, error) {
	var results SearchResults

	results.Total = r.Hits.Total.Value

	// handle empty results
	if len(r.Hits.Hits) < 1 {
		results.Results = []interface{}{}
		return results, nil
	}

	// unmarshal elasticsearch hits
	for _, hit := range r.Hits.Hits {
		h, err := doc.UnmarshalHit(hit)
		if err != nil {
			return results, err
		}
		results.Results = append(results.Results, h)
	}

	return results, nil
}

// ReadSearchResponse reads an Elasticsearch response for a Search request
// and returns a SearchResults or the first non-nil error occurring
// in the process.
//
// It must be provided an Unmarshaler to determinate the marshaling process.
// The typical usage is to provide an entity having a custom NewHit method
// (see Document interface).
func ReadSearchResponse(res *esapi.Response, doc Unmarshaler) (SearchResults, error) {
	if res.IsError() {
		return SearchResults{}, statusError(res.StatusCode)
	}

	r, err := decodeRawSearchResponse(res)
	if err != nil {
		return SearchResults{}, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	return r.unwrap(doc)
}

// decodeRawSearchResponse decodes a raw json response from Elasticsearch
// as an esSearchResponse.
func decodeRawSearchResponse(res *esapi.Response) (esSearchResponse, error) {
	var r esSearchResponse
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return r, err
	}
	return r, nil
}
