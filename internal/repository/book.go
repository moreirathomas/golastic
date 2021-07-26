package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/moreirathomas/golastic/internal"
	"github.com/moreirathomas/golastic/pkg/golastic"
)

// Ensure Repository implements BookService
var _ internal.BookService = (*Repository)(nil)

// SearchBooks retrieves books matching the userQuery in the database
// or the first non-nil error encountered in the process.
func (r Repository) SearchBooks(userQuery string, size, from int) ([]internal.Book, int, error) {
	esQuery := buildSearchQuery(userQuery, size, from)
	res, err := r.makeSearch(esQuery)
	if err != nil {
		return []internal.Book{}, 0, err
	}

	books, err := unmarshalHits(res.Hits)
	if err != nil {
		return books, 0, fmt.Errorf("failed to unmarshal books: %w", err)
	}

	return books, res.Total, nil
}

// makeSearch performs an Elasticsearch search request with the given query.
func (r *Repository) makeSearch(esQuery io.Reader) (golastic.SearchResults, error) {
	res, err := r.es.Search(
		r.es.Search.WithIndex(r.indexName),
		r.es.Search.WithBody(esQuery),
		r.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return golastic.SearchResults{}, err
	}

	defer res.Body.Close()
	if err := golastic.ReadErrorResponse(res); err != nil {
		return golastic.SearchResults{}, err
	}

	return golastic.UnwrapSearchResponse(res, internal.Book{})
}

// buildSearchQuery builds an Elasticsearch search query.
func buildSearchQuery(s string, size, from int) io.Reader {
	if s == "" {
		return golastic.MatchAllSearchQuery(size, from).Reader()
	}

	q := golastic.NewSearchQuery(s, golastic.SearchQueryConfig{
		Fields: []golastic.Field{
			{Name: "title", Weight: 10},
			{Name: "abstract"},
		},
		Sort: []map[string]string{
			{"_score": "asc"},
			{"_doc": "asc"},
		},
		Size: size,
		From: from,
	})

	return q.Reader()
}

func unmarshalHits(hits []interface{}) ([]internal.Book, error) {
	books := make([]internal.Book, 0, len(hits))
	for _, h := range hits {
		b, ok := h.(internal.Book)
		if !ok {
			return books, fmt.Errorf("hit has invalid book format: %#v", h)
		}
		books = append(books, b)
	}
	return books, nil
}

func (r Repository) GetBookByID(id string) (internal.Book, error) {
	res, err := r.makeGet(id)
	if err != nil {
		return internal.Book{}, err
	}

	book, ok := res.(internal.Book)
	if !ok {
		return book, fmt.Errorf("response has invalid book format: %#v", res)
	}

	return book, nil
}

// makeGet performs an Elasticsearch get request with the given id.
func (r Repository) makeGet(id string) (interface{}, error) {
	res, err := r.es.Get(r.indexName, id)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if err := golastic.ReadErrorResponse(res); err != nil {
		return nil, err
	}

	return golastic.UnwrapGetResponse(res, internal.Book{})
}

// InsertBook indexes a new book.
func (r Repository) InsertBook(b internal.Book) error {
	payload, err := json.Marshal(b)
	if err != nil {
		return err
	}

	res, err := r.es.Index(r.indexName, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return golastic.ReadErrorResponse(res)
}

// InsertManyBooks indexes multiple new book documents at once.
func (r *Repository) InsertManyBooks(books []internal.Book) error {
	in := make([]interface{}, len(books))
	for i, b := range books {
		in[i] = b
	}

	cfg := golastic.ContextConfig{
		IndexName: r.indexName,
		Client:    r.es,
	}

	return golastic.BulkIndex(cfg, in)
}

// UpdateBook updates the specified book with a partial book input.
func (r Repository) UpdateBook(b internal.Book) error {
	// The document must be wrapped in a "doc" object
	payload, err := json.Marshal(map[string]internal.Book{
		"doc": b,
	})
	if err != nil {
		return err
	}

	res, err := r.es.Update(r.indexName, b.ID, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return golastic.ReadErrorResponse(res)
}

// DeleteBook removes the specified book from the index.
func (r Repository) DeleteBook(id string) error {
	res, err := r.es.Delete(r.indexName, id)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return golastic.ReadErrorResponse(res)
}
