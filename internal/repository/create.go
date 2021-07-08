package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esutil"
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

// CreateBulk indexes multiple new book documents at once.
func (r *Repository) CreateBulk(books []internal.Book) error {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         r.indexName,
		Client:        r.es,
		NumWorkers:    runtime.NumCPU(),
		FlushBytes:    int(5e+6),
		FlushInterval: 30 * time.Second,
	})
	if err != nil {
		return err
	}

	for _, v := range books {
		var err error
		payload, err := json.Marshal(v)
		if err != nil {
			return err
		}

		err = bi.Add(context.Background(), esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewReader(payload),
			OnFailure: func(c context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem, e error) {
				err = fmt.Errorf("error: %s", e)
			},
		})

		if err != nil {
			return err
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		return err
	}

	biStats := bi.Stats()

	if biStats.NumFailed > 0 {
		log.Printf("indexed [%d] documents with [%d] errors", biStats.NumFlushed, biStats.NumFailed)
	} else {
		log.Printf("Sucessfuly indexed [%d] documents", biStats.NumFlushed)
	}

	return nil
}
