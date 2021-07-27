package golastic

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type SearchAPI struct {
	client *elasticsearch.Client
	index  string
}

func (api SearchAPI) MatchAllQuery(p SearchPagination) (*esapi.Response, error) {
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

	return res, nil
}

func (api SearchAPI) MultiMatchQuery(qs string, f []Field, p SearchPagination, s SearchSort) (*esapi.Response, error) {
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

	return res, nil
}
