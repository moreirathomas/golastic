package golastic

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type SearchAPI struct {
	ContextConfig
}

func (api SearchAPI) MatchAllQuery(size, from int) (*esapi.Response, error) {
	// TODO
	q := MatchAllSearchQuery(size, from)

	raw, err := api.Client.Search(
		api.Client.Search.WithIndex(api.IndexName),
		api.Client.Search.WithBody(q.Reader()),
		api.Client.Search.WithTrackTotalHits(true),
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
	q := NewSearchQuery(qs, cfg)

	raw, err := api.Client.Search(
		api.Client.Search.WithIndex(api.IndexName),
		api.Client.Search.WithBody(q.Reader()),
		api.Client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: failed to perform search: %s",
			ErrBadRequest, err,
		)
	}

	return raw, nil
}
