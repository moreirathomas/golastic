package golastic

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type SearchAPI struct {
	client *elasticsearch.Client
	index  string
}

func (api *SearchAPI) MatchAllQuery(p SearchPagination) (*SearchResults, error) {
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

func (api *SearchAPI) MultiMatchQuery(qs string, f []Field, p SearchPagination, s SearchSort) (*SearchResults, error) {
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

type SearchResults struct {
	Hits *SearchHits `json:"hits,omitempty"` // the actual search hits
}

func (r *SearchResults) TotalHits() int {
	if r != nil && r.Hits != nil && r.Hits.Total != nil {
		return r.Hits.Total.Value
	}
	return 0
}

func (r *SearchResults) UnwrapHits(doc Unmarshaler) ([]interface{}, error) {
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

type SearchHits struct {
	Total *TotalHits `json:"total,omitempty"` // total number of hits found
	Hits  []*Hit     `json:"hits,omitempty"`  // the actual hits returned
}

type TotalHits struct {
	Value int `json:"value"` // value of the total hit count
}

func decodeSearchResults(res *esapi.Response) (*SearchResults, error) {
	defer res.Body.Close()
	if err := ReadErrorResponse(res); err != nil {
		return nil, err
	}

	var r SearchResults
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
