package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/moreirathomas/golastic/internal"
)

// Create indexes a new document.
func (r *Repository) Create(doc internal.Book) error {
	payload, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	ctx := context.Background()

	req := esapi.CreateRequest{
		Index:      r.indexName,
		DocumentID: strconv.Itoa(doc.ID), // FIXME the ID must be known before adding to ES!
		Body:       bytes.NewReader(payload),
	}

	res, err := req.Do(ctx, r.es)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}

	return nil
}
