package golastic

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/clarketm/json"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

type DocumentAPI struct {
	ContextConfig
}

func (api DocumentAPI) Get(id string) (*esapi.Response, error) {
	res, err := api.Client.Get(api.IndexName, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}
	return res, nil
}

func (api DocumentAPI) Update(id string, doc interface{}) (*esapi.Response, error) {
	// Elasticsearch expects the document to be wrapped inside
	// an object with "doc" key.
	payload, err := json.Marshal(map[string]interface{}{"doc": doc})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	res, err := api.Client.Update(api.IndexName, id, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	return res, nil
}

func (api DocumentAPI) Index(doc interface{}) (*esapi.Response, error) {
	payload, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	res, err := api.Client.Index(api.IndexName, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	return res, nil
}

func (api DocumentAPI) Delete(id string) (*esapi.Response, error) {
	res, err := api.Client.Delete(api.IndexName, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	return res, nil
}

func (api DocumentAPI) Bulk(docs []interface{}) error {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  api.IndexName,
		Client: api.Client,
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

	stats := bi.Stats()
	if stats.NumFailed > 0 {
		log.Printf("indexed [%d] documents with [%d] errors", stats.NumFlushed, stats.NumFailed)
	} else {
		log.Printf("Successfully indexed [%d] documents", stats.NumFlushed)
	}

	return nil
}
