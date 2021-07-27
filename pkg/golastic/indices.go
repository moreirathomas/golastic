package golastic

import (
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

type IndicesAPI struct {
	client *elasticsearch.Client
}

// Exists returns true when the index already exists.
func (api IndicesAPI) Exists(index string) (bool, error) {
	res, err := api.client.Indices.Exists([]string{index})
	if err != nil {
		return false, err
	}
	switch res.StatusCode {
	case 200:
		return true, nil
	case 404:
		return false, nil
	default:
		return false, fmt.Errorf("[%s]", res.Status())
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

	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}

	return nil
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
