package http

import (
	"net/http"
	"time"

	"github.com/moreirathomas/golastic/internal"
	"github.com/moreirathomas/golastic/pkg/golastic"
	"github.com/moreirathomas/golastic/pkg/httputil"
)

// SearchBooks retrieves all books matching the query string,
// either in their title or in their abstract.
func (s Server) SearchBooks(w http.ResponseWriter, r *http.Request) {
	// Retrieve user's query string
	q := extractQueryParam(r, "query")

	// Retrieve pagination parameters
	size, err := extractQueryParamAsInt(r, "size")
	if err != nil {
		// FIXME move this default declaration to internal/repository
		// Elasticsearch should be fine if we omit these params (do tests though)
		size = golastic.DefaultQuerySize
	}
	from, err := extractQueryParamAsInt(r, "from")
	if err != nil {
		// FIXME
		from = golastic.DefaultQueryFrom
	}

	// Perform ElasticSearch query
	results, total, err := s.Repository.SearchBooks(q, size, from)
	if err != nil {
		respondHTTPError(w, errInternal.Wrap(err))
		return
	}

	res := struct {
		Data interface{} `json:"data"`
		httputil.Pagination
	}{
		Data: struct {
			Results []internal.Book `json:"results"`
			Total   int             `json:"total"`
		}{Results: results, Total: total},
		Pagination: httputil.NewPagination(size, from),
	}
	res.Pagination.SetLinks(r)

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
	if err := s.Repository.InsertBook(book); err != nil {
		// TODO: specify error handling (could be a duplicate or internal error)
		respondHTTPError(w, errBadRequest.Wrap(err))
		return
	}

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
