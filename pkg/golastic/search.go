// This file regroups all entities and methods to interact with
// Elasticseach Search API.

package golastic

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// SearchAPI is used to search for documents in Elasticsearch.
type SearchAPI struct {
	client *elasticsearch.Client
	index  string
}

// MatchAllQuery returns the result of a query which match all documents.
func (api *SearchAPI) MatchAllQuery(p SearchPagination) (*SearchResult, error) {
	res, err := api.client.Search(
		api.client.Search.WithIndex(api.index),
		api.client.Search.WithBody(newMatchAllQuery().Reader()),
		api.client.Search.WithSort(defaultSort...),
		api.client.Search.WithFrom(p.From),
		api.client.Search.WithSize(p.Size),
		api.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to perform search: %s", ErrBadRequest, err)
	}

	r, err := decodeSearchResults(res)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// MultiMatchQuery returns the result of a query which performs
// a full text query across multiple fields.
func (api *SearchAPI) MultiMatchQuery(qs string, f []Field, p SearchPagination, s SearchSort) (*SearchResult, error) {
	if len(s) == 0 {
		s = defaultSort
	}

	res, err := api.client.Search(
		api.client.Search.WithIndex(api.index),
		api.client.Search.WithBody(newMultiMatchQuery(qs, f).Reader()),
		api.client.Search.WithSort(s...),
		api.client.Search.WithFrom(p.From),
		api.client.Search.WithSize(p.Size),
		api.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to perform search: %s", ErrBadRequest, err)
	}

	r, err := decodeSearchResults(res)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// SearchResult is the result of search in Elasticsearch.
type SearchResult struct {
	Hits *SearchHits `json:"hits,omitempty"`
}

// TotalHits conveniently returns the number of hits for a search result.
func (r *SearchResult) TotalHits() int {
	if r != nil && r.Hits != nil && r.Hits.Total != nil {
		return r.Hits.Total.Value
	}
	return 0
}

// UnwrapHits conveniently returns the response hits. Each hit is unmarshalled
// based on the given Unmarshaler parameter and returned as an interface left
// to be type asserted by the caller.
func (r *SearchResult) UnwrapHits(doc Unmarshaler) ([]interface{}, error) {
	if r.Hits == nil || r.Hits.Hits == nil || len(r.Hits.Hits) == 0 {
		return nil, nil
	}

	hits := make([]interface{}, 0, len(r.Hits.Hits))
	for _, hit := range r.Hits.Hits {
		h, err := doc.UnmarshalHit(*hit)
		if err != nil {
			return nil, err
		}
		hits = append(hits, h)
	}

	return hits, nil
}

// SearchHits represents the definition of the list of hits.
type SearchHits struct {
	Total *TotalHits `json:"total,omitempty"`
	Hits  []*Hit     `json:"hits,omitempty"` // The actual hits returned.
}

// SearchHits is the total number of hits.
type TotalHits struct {
	Value int `json:"value"`
}

func decodeSearchResults(res *esapi.Response) (*SearchResult, error) {
	defer res.Body.Close()
	if err := ReadErrorResponse(res); err != nil {
		return nil, err
	}

	var r SearchResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
