package http

import (
	"net/http"
	"time"

	"github.com/moreirathomas/golastic/pkg/golastic"
	"github.com/moreirathomas/golastic/pkg/pagination"
)

// SearchBooks retrieves all books matching the query string,
// either in their title or in their abstract.
func (s Server) SearchBooks(w http.ResponseWriter, r *http.Request) {
	// Retrieve user's query string
	q := extractQueryParam(r, "query")

	// Retrieve pagination parameters
	size, err := extractQueryParamInt(r, "size")
	if err != nil {
		size = golastic.DefaultQuerySize
	}
	page, err := extractQueryParamInt(r, "page")
	if err != nil || page < 1 {
		page = 1
	}
	from := pagination.PageToOffset(page, size)

	// Perform ElasticSearch query
	results, total, err := s.Repository.SearchBooks(q, size, from)
	if err != nil {
		respondHTTPError(w, errInternal.Wrap(err))
		return
	}

	// Paginate the results and send the response
	p, err := pagination.New(r, total, page, size)
	if err != nil {
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	res := struct {
		Results interface{} `json:"results"`
		Total   int         `json:"total"`
		pagination.Pagination
	}{
		Results:    results,
		Total:      total,
		Pagination: p,
	}

	respondJSON(w, 200, res)
}

// GetBookByID retrieves a book by its ID in the repository.
func (s Server) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractRouteParam(r, "bookID")
	if err != nil {
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	book, err := s.Repository.GetBookByID(id)
	if err != nil {
		respondHTTPError(w, errNotFound.Wrap(err))
		return
	}

	respondJSON(w, 200, book)
}

// InsertBook adds a new book in the repository, if the request is valid.
func (s Server) InsertBook(w http.ResponseWriter, r *http.Request) {
	book, err := readBookPayload(r.Body)
	if err != nil {
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	book.CreatedAt = time.Now()
	id, err := s.Repository.InsertBook(book)
	if err != nil {
		// TODO: specify error handling (could be a duplicate or internal error)
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	// Populate the book instance with the ID created on Elasticsearch part.
	book.ID = id

	respondJSON(w, 201, book)
}

// UpdateBook adds a new book in the repository, if the request is valid.
func (s Server) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := extractRouteParam(r, "bookID")
	if err != nil {
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	book, err := readBookPayload(r.Body)
	if err != nil {
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}
	book.ID = id

	if err := s.Repository.UpdateBook(book); err != nil {
		// TODO: specify error handling (could be a duplicate or internal error)
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	respondJSON(w, 204, nil)
}

func (s Server) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := extractRouteParam(r, "bookID")
	if err != nil {
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	if err := s.Repository.DeleteBook(id); err != nil {
		// TODO: specify error handling (could be internal)
		respondHTTPError(w, errNotFound.Wrap(err))
	}

	respondJSON(w, 204, nil)
}
