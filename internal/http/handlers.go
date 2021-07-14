package http

import (
	"net/http"
)

// SearchBooks retrieves all books matching the query string,
// either in their title or in their abstract.
func (s Server) SearchBooks(w http.ResponseWriter, r *http.Request) {
	// Retrieve user's query string
	q, err := s.readURLQuery(r, "query")
	if err != nil {
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	// Perform ElasticSearch query
	results, err := s.Repository.SearchBooks(q)
	if err != nil {
		respondHTTPError(w, errInternal.Wrap(err))
		return
	}

	respondJSON(w, 200, results)
}

// GetBookByID retrieves a book by its ID in the repository.
func (s Server) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "bookID")
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

	if err := s.Repository.InsertBook(book); err != nil {
		// TODO: specify error handling (could be a duplicate or internal error)
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	respondJSON(w, 201, book)
}

// UpdateBook adds a new book in the repository, if the request is valid.
func (s Server) UpdateBook(w http.ResponseWriter, r *http.Request) {
	book, err := readBookPayload(r.Body)
	if err != nil {
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	if err := s.Repository.UpdateBook(book); err != nil {
		// TODO: specify error handling (could be a duplicate or internal error)
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

	respondJSON(w, 204, nil)
}

func (s Server) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "bookID")
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
