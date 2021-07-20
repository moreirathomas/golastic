package golastic

import (
	"bytes"
	"context"
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
	defer bi.Close(context.Background())

	for _, doc := range docs {
		payload, err := json.Marshal(doc)
		if err != nil {
			return err
		}

		if err := bi.Add(context.Background(), esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewReader(payload),
			OnFailure: func(_ context.Context, _ esutil.BulkIndexerItem, _ esutil.BulkIndexerResponseItem, e error) {
				log.Printf("failed to index document %#v: %s", doc, e)
			},
		}); err != nil {
			return err
		}
	}

	logBulkIndexerStats(bi)
	return nil
}

func logBulkIndexerStats(bi esutil.BulkIndexer) {
	stats := bi.Stats()
	if stats.NumFailed > 0 {
		log.Printf("indexed [%d] documents with [%d] errors", stats.NumFlushed, stats.NumFailed)
	} else {
		log.Printf("Successfully indexed [%d] documents", stats.NumFlushed)
	}
}
