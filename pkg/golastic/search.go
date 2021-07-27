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

func (api SearchAPI) MatchAllQuery(size, from int) (*esapi.Response, error) {
	// TODO
	q := matchAllSearchQuery(size, from)

	raw, err := api.client.Search(
		api.client.Search.WithIndex(api.index),
		api.client.Search.WithBody(q.Reader()),
		api.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: failed to perform search: %s",
			ErrBadRequest, err,
		)
	}

	return raw, nil
}

func (api SearchAPI) MultiMatchQuery(qs string, cfg SearchQueryConfig) (*esapi.Response, error) {
	// TODO
	q := newSearchQuery(qs, cfg)

	raw, err := api.client.Search(
		api.client.Search.WithIndex(api.index),
		api.client.Search.WithBody(q.Reader()),
		api.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: failed to perform search: %s",
			ErrBadRequest, err,
		)
	}

	return raw, nil
}
