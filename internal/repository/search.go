package repository

import (
	"fmt"
	"io"

	"github.com/moreirathomas/golastic/internal"
	"github.com/moreirathomas/golastic/pkg/golastic"
)

// Search returns results matching a query.
func (r *Repository) Search(query string) (*golastic.SearchResults, error) {
	var results golastic.SearchResults

	res, err := r.es.Search(
		r.es.Search.WithIndex(r.indexName),
		r.es.Search.WithBody(buildSearchQuery(query)),
		r.es.Search.WithTrackTotalHits(true),
		r.es.Search.WithPretty(),
	)
	if err != nil {
		return &results, err
	}

	defer res.Body.Close()
	if res.IsError() {
		return &results, fmt.Errorf("error: %s", res)
	}

	results, err = golastic.UnwrapSearchResponse(res, internal.Book{})
	if err != nil {
		return &results, err
	}

	return &results, nil
}

func buildSearchQuery(search string) io.Reader {
	if search == "" {
		return golastic.MatchAllSearchQuery(10).Reader()
	}
	return golastic.NewSearchQuery(search, golastic.SearchQueryConfig{
		Fields: []golastic.Field{
			{Name: "title", Weight: 10},
			{Name: "abstract"},
		},
		Sort: []map[string]string{
			{"_score": "asc"},
			{"_doc": "asc"},
		},
		Size: 25,
	}).Reader()
}