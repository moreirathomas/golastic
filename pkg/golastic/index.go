package golastic

import (
	"fmt"
	"strings"
)

// IndexExists returns true when the index already exists in the repository.
func IndexExists(c ContextConfig) (bool, error) {
	res, err := c.Client.Indices.Exists([]string{c.IndexName})
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

// CreateIndex creates a new index with mapping.
func CreateIndex(c ContextConfig, mapping string) error {
	res, err := c.Client.Indices.Create(
		c.IndexName,
		c.Client.Indices.Create.WithBody(strings.NewReader(mapping)),
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

// CreateIndexIfNotExists creates a new index with mapping
// if the index does not exists yet on the client.
// It returns true if the index is being created.
func CreateIndexIfNotExists(c ContextConfig, mapping string) (bool, error) {
	exists, err := IndexExists(c)
	switch {
	case err != nil:
		return false, err
	case exists:
		return false, nil
	default:
		return true, CreateIndex(c, mapping)
	}
}