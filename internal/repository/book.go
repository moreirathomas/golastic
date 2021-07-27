package repository

import (
	"fmt"

	"github.com/moreirathomas/golastic/internal"
	"github.com/moreirathomas/golastic/pkg/golastic"
)

// Ensure Repository implements BookService
var _ internal.BookService = (*Repository)(nil)

// SearchBooks retrieves books matching the userQuery in the database
// or the first non-nil error encountered in the process.
func (r Repository) SearchBooks(userQuery string, size, from int) ([]internal.Book, int, error) {
	handleError := func(err error) ([]internal.Book, int, error) {
		return []internal.Book{}, 0, err
	}

	esQuery := buildSearchQuery(userQuery, size, from)

	res, err := esQuery.Do(r.esContext)
	if err != nil {
		return handleError(err)
	}

	results, err := golastic.ReadSearchResponse(res, internal.Book{})
	if err != nil {
		return handleError(err)
	}

	books, err := unmarshalHits(results.Results)
	if err != nil {
		return handleError(fmt.Errorf("failed to unmarshal books: %w", err))
	}

	return books, results.Total, nil
}

// buildSearchQuery builds an Elasticsearch search query.
func buildSearchQuery(s string, size, from int) golastic.SearchQuery {
	if s == "" {
		return golastic.MatchAllSearchQuery(size, from)
	}

	return golastic.NewSearchQuery(s, golastic.SearchQueryConfig{
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
	res, err := golastic.Get(r.esContext, id)
	if err != nil {
		return internal.Book{}, err
	}

	result, err := golastic.ReadGetResponse(res, internal.Book{})
	if err != nil {
		return internal.Book{}, err
	}

	book, ok := result.(internal.Book)
	if !ok {
		return book, fmt.Errorf("response has invalid book format: %#v", res)
	}

	return book, nil
}

// InsertBook indexes a new book.
func (r Repository) InsertBook(b internal.Book) (string, error) {
	res, err := golastic.Insert(r.esContext, b)
	if err != nil {
		return "", fmt.Errorf(
			"%w failed to insert book %#v: %s",
			ErrInternal, b, err,
		)
	}

	id, err := golastic.ReadInsertResponse(res)
	if err != nil {
		return "", fmt.Errorf("could not insert book: %w", err)
	}

	return id, nil
}

// InsertManyBooks indexes multiple new book documents at once.
func (r *Repository) InsertManyBooks(books []internal.Book) error {
	in := make([]interface{}, len(books))
	for i, b := range books {
		in[i] = b
	}

	if err := golastic.BulkIndex(r.esContext, in); err != nil {
		return fmt.Errorf(
			"%w: failed to insert books: %s",
			ErrInternal, err,
		)
	}

	return nil
}

// UpdateBook updates the specified book with a partial book input.
func (r Repository) UpdateBook(b internal.Book) error {
	res, err := golastic.Update(r.esContext, b.ID, b)
	if err != nil {
		return fmt.Errorf(
			"%w: failed to update book %#v: %s",
			ErrInternal, b, err,
		)
	}

	// TODO: handle 404 when golastic allows it
	return golastic.ReadErrorResponse(res)
}

// DeleteBook removes the specified book from the index.
func (r Repository) DeleteBook(id string) error {
	res, err := golastic.Delete(r.esContext, id)
	if err != nil {
		return err
	}

	// TODO: handle 404 when golastic allows it
	return golastic.ReadErrorResponse(res)
}
