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
