package golastic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/clarketm/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

type DocumentAPI struct {
	client *elasticsearch.Client
	index  string
}

// -- Get API

func (api *DocumentAPI) Get(id string) (*GetResult, error) {
	res, err := api.client.Get(api.index, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	defer res.Body.Close()
	if err := ReadErrorResponse(res); err != nil {
		return nil, err
	}

	var r GetResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

type GetResult struct {
	Found bool `json:"found"`
	Hit
}

func (r *GetResult) Unwrap(doc Unmarshaler) (interface{}, error) {
	if !r.Found {
		return nil, errors.New("not found") // TODO is it really an error?
	}

	result, err := doc.UnmarshalHit(r.Hit)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// -- Update API

func (api *DocumentAPI) Update(id string, doc interface{}) (*esapi.Response, error) {
	// Elasticsearch expects the document to be wrapped inside
	// an object with "doc" key.
	payload, err := json.Marshal(map[string]interface{}{"doc": doc})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	res, err := api.client.Update(api.index, id, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	return res, nil
}

// -- Index API

func (api *DocumentAPI) Index(doc interface{}) (*IndexResult, error) {
	payload, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	res, err := api.client.Index(api.index, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	defer res.Body.Close()
	if err := ReadErrorResponse(res); err != nil {
		return nil, err
	}

	var r IndexResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

type IndexResult struct {
	ID     string `json:"_id"`
	Result string `json:"result"` // "created" in case of success
}

// TODO This may be useless. In which scenario we don't get an error
// and yet the doc is not created ?
func (r *IndexResult) Unwrap() (string, error) {
	if r.Result != "created" {
		return "", errors.New("not created")
	}

	return r.ID, nil
}

// -- Delete API

func (api *DocumentAPI) Delete(id string) (*esapi.Response, error) {
	res, err := api.client.Delete(api.index, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	return res, nil
}

// -- Bulk API

func (api *DocumentAPI) Bulk(docs []interface{}) error {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  api.index,
		Client: api.client,
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
