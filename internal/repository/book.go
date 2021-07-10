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
func (r Repository) SearchBooks(userQuery string) ([]internal.Book, error) {
	res, err := r.search(userQuery)
	if err != nil {
		return []internal.Book{}, err
	}

	var books []internal.Book
	for _, hit := range res.Hits {
		b, ok := hit.(internal.Book)
		if !ok {
			return books, fmt.Errorf("hit has invalid book format: %#v", hit)
		}
		books = append(books, b)
	}

	return books, nil
}

// search is a helper to encapsulate Elasticsearch interface call.
func (r *Repository) search(query string) (*golastic.SearchResults, error) {
	var results golastic.SearchResults

	res, err := r.es.Search(
		r.es.Search.WithIndex(r.indexName),
		r.es.Search.WithBody(buildSearchQuery(query)),
		r.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return &results, err
	}

	defer res.Body.Close()
	if err := golastic.ReadErrorResponse(res); err != nil {
		return &results, err
	}

	results, err = golastic.UnwrapSearchResponse(res, internal.Book{})
	if err != nil {
		return &results, err
	}

	return &results, nil
}

// buildSearchQuery is a helper to encapsulate Elasticsearch interface call.
func buildSearchQuery(s string) io.Reader {
	if s == "" {
		return golastic.MatchAllSearchQuery(10).Reader()
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
		Size: 25,
	})

	return q.Reader()
}

func (r Repository) GetBookByID(id string) (internal.Book, error) {
	// FIXME we need to declare a type for GET api before
	// unmarshalling the response to a Book.
	// return r.Get(id)
	return internal.Book{}, nil
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
	if err := golastic.ReadErrorResponse(res); err != nil {
		return err
	}

	return nil
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
	payload, err := json.Marshal(b)
	if err != nil {
		return err
	}

	res, err := r.es.Update(r.indexName, b.ID, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if err := golastic.ReadErrorResponse(res); err != nil {
		return err
	}

	return nil
}

// DeleteBook removes the specified book from the index.
func (r Repository) DeleteBook(id string) error {
	res, err := r.es.Delete(r.indexName, id)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if err := golastic.ReadErrorResponse(res); err != nil {
		return err
	}

	return nil
}
