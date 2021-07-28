// This file regroups all entities and methods to interact with
// Elasticseach single document APIs, namely Index, Get, Delete
// and Update APIs.
// It also contains one entity to interact with the  multi document
// Bulk API.

package golastic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/clarketm/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

// DocumentAPI is used to interact with documents in Elasticsearch.
type DocumentAPI struct {
	client *elasticsearch.Client
	index  string
}

// -- Get API

// Get returns the result of a getting a document in Elasticsearch.
func (api *DocumentAPI) Get(id string) (*GetResult, error) {
	res, err := api.client.Get(api.index, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	defer res.Body.Close()
	if err := readErrorResponse(res); err != nil {
		return nil, err
	}

	var r GetResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

// IndexResult is the result of getting a document in Elasticsearch.
type GetResult struct {
	Found bool `json:"found"`
	Hit
}

// Unwrap conveniently returns the response hit. The hit is unmarshalled
// based on the given Unmarshaler parameter and returned as an interface left
// to be type asserted by the caller.
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

// Update returns the result of updating a document in Elasticsearch.
func (api *DocumentAPI) Update(id string, doc interface{}) error {
	// Elasticsearch expects the document to be wrapped inside
	// an object with "doc" key.
	payload, err := json.Marshal(map[string]interface{}{"doc": doc})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	res, err := api.client.Update(api.index, id, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	defer res.Body.Close()
	return readErrorResponse(res)
}

// -- Index API

// Update returns the result of a indexing a document in Elasticsearch.
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
	if err := readErrorResponse(res); err != nil {
		return nil, err
	}

	var r IndexResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

// IndexResult is the result of indexing a document in Elasticsearch.

type IndexResult struct {
	ID     string `json:"_id"`
	Result string `json:"result"` // "created" in case of success
}

// TODO This may be useless. In which scenario we don't get an error
// and yet the doc is not created ?
// Unwrap conveniently returns the document ID.
func (r *IndexResult) Unwrap() (string, error) {
	if r.Result != "created" {
		return "", errors.New("not created")
	}

	return r.ID, nil
}

// -- Delete API

// Update returns the result of a deleting a document in Elasticsearch.
func (api *DocumentAPI) Delete(id string) error {
	res, err := api.client.Delete(api.index, id)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnhandled, err)
	}

	defer res.Body.Close()
	return readErrorResponse(res)
}

// -- Bulk API

// Update returns the result of a indexing many documents in Elasticsearch.
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
