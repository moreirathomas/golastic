// This file regroups all entities and methods to interact with
// Elasticseach Indices API.

package golastic

import (
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

// SearchAPI is used to interact with indices in Elasticsearch.
type IndicesAPI struct {
	client *elasticsearch.Client
}

// Exists returns true when the index already exists.
func (api IndicesAPI) Exists(index string) (bool, error) {
	res, err := api.client.Indices.Exists([]string{index})
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

// Create creates a new index with mapping.
func (api IndicesAPI) Create(index, mapping string) error {
	res, err := api.client.Indices.Create(
		index,
		api.client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return err
	}

	return ReadErrorResponse(res)
}

// CreateIfNotExists creates a new index with mapping if the index does not
// exists on the client. It returns true if the index is being created.
func (api IndicesAPI) CreateIfNotExists(index, mapping string) (bool, error) {
	exists, err := api.Exists(index)
	switch {
	case err != nil:
		return false, err
	case exists:
		return false, nil
	default:
		return true, api.Create(index, mapping)
	}
}
