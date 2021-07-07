package internal

import "time"

// Book represents a book in the API.
type Book struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Abstract  string    `json:"abstract"`
	Author    Author    `json:"author"`
}

// Author represents a book's author.
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// BookService gathers repository methods to perform CRUD on books.
type BookService interface {

	// SearchBooks retrieves all books matching the input query.
	SearchBooks(query string) ([]Book, error)

	// GetBookByID retrieves a book by its ID in the repository.
	// It returns a non-nil error if one occurs in the process
	// or if no match were found.
	GetBookByID(id int) (Book, error)

	// InsertBook adds the given book in the repository.
	InsertBook(book Book) error

	// UpdateBook updates a book in the repository.
	UpdateBook(book Book) error

	// DeleteBook deletes a book by its ID in the repository.
	DeleteBook(id int) error
}

// Validate is a WIP. Its implementation remains to be done.
//
// Validate return a non-nil error if the book receiver does not match
// the validation requirements.
func (b Book) Validate() error {
	return nil
}
