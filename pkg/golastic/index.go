package golastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// IndexExists returns true when the index already exists in the repository.
func IndexExists(ctx ContextConfig) (bool, error) {
	res, err := ctx.Client.Indices.Exists([]string{ctx.IndexName})
	if err != nil {
		return false, err
	}
	switch err := ReadErrorResponse(res); err {
	case nil:
		return true, nil
	case ErrNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("[%s] %w", res.Status(), err)
	}
}

// CreateIndex creates a new index with mapping.
func CreateIndex(ctx ContextConfig, mapping string) error {
	res, err := ctx.Client.Indices.Create(
		ctx.IndexName,
		ctx.Client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return err
	}

	return ReadErrorResponse(res)
}

// CreateIndexIfNotExists creates a new index with mapping
// if the index does not exists yet on the client.
// It returns true if the index is being created.
func CreateIndexIfNotExists(ctx ContextConfig, mapping string) (bool, error) {
	exists, err := IndexExists(ctx)
	switch {
	case err != nil:
		return false, err
	case exists:
		return false, nil
	default:
		return true, CreateIndex(ctx, mapping)
	}
}

// getResponseWrapper represents selected fields from
// the response to an Elasticsearch Document Indexing request.
type indexResponseWrapper struct {
	Result string `json:"result"`
	ID     string `json:"_id"`
}

// UnwrapIndexResponse reads an Elasticsearch response for a Document Indexing
// request and returns a string corresponding to the ID of the indexed document
// or the first non-nil error occurring in the process.
func UnwrapIndexResponse(res *esapi.Response) (string, error) {
	var rw indexResponseWrapper
	if err := json.NewDecoder(res.Body).Decode(&rw); err != nil {
		return "", err
	}

	if rw.Result != "created" {
		return "", errors.New("not created")
	}

	return rw.ID, nil
}
