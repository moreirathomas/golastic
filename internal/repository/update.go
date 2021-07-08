package repository

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/moreirathomas/golastic/internal"
)

func (r *Repository) Update(book internal.Book) error {
	payload, err := json.Marshal(book)
	if err != nil {
		return err
	}

	res, err := r.es.Update(r.indexName, book.ID, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}

	return nil
}
