package repository

import (
	"log"

	"github.com/moreirathomas/golastic/internal"
)

// Ensure Repository implements BookService
var _ internal.BookService = (*Repository)(nil)

// SearchBooks is a WIP.
// It remains to be implemented and its signature might change.
// It will likely return and ES result type in the future.
func (r Repository) SearchBooks(q string) ([]internal.Book, error) {
	res, err := r.Search(q)
	if err != nil {
		return []internal.Book{}, err
	}

	log.Printf("Retrieved %d books\n", res.Total)
	// TODO Hit may not embed Book in the future.
	var books []internal.Book
	for _, hit := range res.Hits {
		b := internal.Book{
			ID:        hit.ID,
			CreatedAt: hit.CreatedAt,
			Title:     hit.Title,
			Abstract:  hit.Abstract,
			Author:    hit.Author,
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
