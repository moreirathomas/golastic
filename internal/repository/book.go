package repository

import (
	"fmt"
	"log"

	"github.com/moreirathomas/golastic/internal"
)

// Ensure Repository implements BookService
var _ internal.BookService = (*Repository)(nil)

// SearchBooks retrieves books matching the userQuery in the database
// or the first non-nil error encountered in the process.
func (r Repository) SearchBooks(userQuery string) ([]internal.Book, error) {
	res, err := r.Search(userQuery)
	if err != nil {
		return []internal.Book{}, err
	}

	log.Printf("Retrieved %d books\n", res.Total)

	books, err := unmarshalBooks(res.Hits)
	if err != nil {
		return books, fmt.Errorf("failed to unmarshal books: %w", err)
	}

	return books, nil
}

func unmarshalBooks(hits []interface{}) ([]internal.Book, error) {
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
	// FIXME we need to declare a type for GET api before
	// unmarshalling the response to a Book.
	// return r.Get(id)
	return internal.Book{}, nil
}

// InsertBook indexes a new book.
func (r Repository) InsertBook(b internal.Book) error {
	return r.Create(b)
}

// UpdateBook updates the specified book with a partial book input.
func (r Repository) UpdateBook(b internal.Book) error {
	return r.Update(b)
}

// DeleteBook removes the specified book from the index.
func (r Repository) DeleteBook(id string) error {
	return r.Delete(id)
}
