package repository

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/moreirathomas/golastic/internal"
)

// Create indexes a new book document.
func (r *Repository) Create(book internal.Book) error {
	payload, err := json.Marshal(book)
	if err != nil {
		return err
	}

	res, err := r.es.Index(r.indexName, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}

	return nil
}
