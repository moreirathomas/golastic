package golastic

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/clarketm/json"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

// BulkIndex indexes the given documents array in bulk.
func BulkIndex(c ContextConfig, docs []interface{}) error {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  c.IndexName,
		Client: c.Client,
	})
	if err != nil {
		return err
	}

	for _, b := range docs {
		payload, err := json.Marshal(b)
		if err != nil {
			return err
		}

		if err := bi.Add(context.Background(), esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewReader(payload),
			OnFailure: func(_ context.Context, _ esutil.BulkIndexerItem, _ esutil.BulkIndexerResponseItem, e error) {
				err = fmt.Errorf("error: %s", e)
			},
		}); err != nil {
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
