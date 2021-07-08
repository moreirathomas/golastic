package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/moreirathomas/golastic/internal"
)

// SearchResults wraps the Elasticsearch search response.
type SearchResults struct {
	Total int    `json:"total"`
	Hits  []*Hit `json:"hits"`
}

// Hit wraps the document returned in search response.
type Hit struct {
	internal.Book
}

// Search returns results matching a query.
func (r *Repository) Search(query string) (*SearchResults, error) {
	var results SearchResults

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

	results, err = unwrapResponse(res)
	if err != nil {
		return &results, err
	}

	return &results, nil
}

type responseWrapper struct {
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

func unwrapResponse(res *esapi.Response) (SearchResults, error) {
	var results SearchResults

	var w responseWrapper
	if err := json.NewDecoder(res.Body).Decode(&w); err != nil {
		return results, err
	}

	results.Total = w.Hits.Total.Value
	if len(w.Hits.Hits) < 1 {
		results.Hits = []*Hit{}
		return results, nil
	}

	for _, hit := range w.Hits.Hits {
		var h Hit
		h.ID = hit.ID

		if err := json.Unmarshal(hit.Source, &h); err != nil {
			return results, err
		}

		results.Hits = append(results.Hits, &h)
	}

	return results, nil
}

func buildSearchQuery(query string) io.Reader {
	var b strings.Builder
	b.WriteString("{\n")

	if query == "" {
		b.WriteString(searchAll)
	} else {
		b.WriteString(fmt.Sprintf(searchMatch, query))
	}

	b.WriteString("\n}")

	return strings.NewReader(b.String())
}

const searchAll = `
	"query" : { "match_all" : {} },
	"size" : 25,
	"sort" : { "created_at" : "desc", "_doc" : "asc" }`

// TODO we could also use "combined_fields"
const searchMatch = `
	"query" : {
		"multi_match" : {
			"query" : %q,
			"fields" : ["title^10", "abstract"],
			"operator" : "and"
		}
	},
	"size" : 25,
	"sort" : [ { "_score" : "desc" }, { "_doc" : "asc" } ]`
