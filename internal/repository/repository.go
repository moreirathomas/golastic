package repository

import (
	"log"
	"strconv"
	"time"

	"github.com/moreirathomas/golastic/internal"
)

var _ internal.BookService = (*Repository)(nil)

type Repository struct{}

// SearchBooks is a WIP.
// It remains to be implemented and its signature might change.
// It will likely return and ES result type in the future.
func (r Repository) SearchBooks(q string) ([]internal.Book, error) {
	log.Println("Searching book: " + q)

	return []internal.Book{
		_mockBook(42),
		_mockBook(314),
		_mockBook(1618),
	}, nil
}

func (r Repository) GetBookByID(id int) (internal.Book, error) {
	log.Println("Getting book with id: " + strconv.Itoa(id))

	return _mockBook(id), nil
}

// InsertBook is a WIP.
// It remains to be implemented and its signature might change.
func (r Repository) InsertBook(b internal.Book) error {
	log.Println("Inserting book: " + b.Title)

	return nil
}

// UpdateBook is a WIP.
// It remains to be implemented and its signature might change.
func (r Repository) UpdateBook(b internal.Book) error {
	log.Println("Updating book: " + b.Title)

	return nil
}

// DeleteBook is a WIP.
// It remains to be implemented and its signature might change.
func (r Repository) DeleteBook(id int) error {
	log.Println("Deleting book with id: " + strconv.Itoa(id))

	return nil
}

// _mockBook is a temporary helper for testing purposes.
// It must be deleted as soon as the crud is ready.
func _mockBook(id int) internal.Book {
	return internal.Book{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "The Fellowship of the Ring",
		Abstract:  "Some cool guys go on a trip.",
		Author: internal.Author{
			Firstname: "Jean-Raoul-Roger",
			Lastname:  "Tolkien",
		},
	}
}
