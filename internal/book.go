package internal

import (
	"encoding/json"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Book represents a book in the API.
type Book struct {
	ID        string    `json:"id,omitempty"`
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
	// It also returns the number of retrieved books.
	SearchBooks(query string, size, from int) ([]Book, int, error)

	// GetBookByID retrieves a book by its ID in the repository.
	// It returns a non-nil error if one occurs in the process
	// or if no match were found.
	GetBookByID(id string) (Book, error)

	// InsertBook adds the given book in the repository.
	InsertBook(book Book) error

	// UpdateBook updates a book in the repository.
	UpdateBook(book Book) error

	// DeleteBook deletes a book by its ID in the repository.
	DeleteBook(id string) error
}

// Validate return a non-nil error if the book receiver does not match
// the validation requirements.
func (b Book) Validate(partial bool) error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Title, validation.Required.When(!partial), validation.Length(1, 100)),
		validation.Field(&b.Abstract, validation.Required.When(!partial)),
		validation.Field(&b.Author, validation.By(func(_ interface{}) error {
			return b.Author.Validate(partial)
		})),
	)
}

// Validate return a non-nil error if the author receiver does not match
// the validation requirements.
func (a Author) Validate(partial bool) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Firstname, validation.Required.When(!partial), is.ASCII, validation.Length(1, 50)),
		validation.Field(&a.Lastname, validation.Required.When(!partial), is.ASCII, validation.Length(1, 50)),
	)
}

// NewHit returns a new hit for ElasticSearch search result that can be later
// casted as a Book. It is necessary to implement elasticsearch.Document interface.
func (b Book) NewHit(id string, src json.RawMessage) (interface{}, error) {
	hit := Book{
		ID: id,
	}
	if err := json.Unmarshal(src, &hit); err != nil {
		return Book{}, err
	}

	return hit, nil
}
