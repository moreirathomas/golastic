package golastic

import (
	"fmt"
	"strings"
)

// IndexExists returns true when the index already exists in the repository.
func IndexExists(ctx ContextConfig) (bool, error) {
	res, err := ctx.Client.Indices.Exists([]string{ctx.IndexName})
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
func CreateIndex(ctx ContextConfig, mapping string) error {
	res, err := ctx.Client.Indices.Create(
		ctx.IndexName,
		ctx.Client.Indices.Create.WithBody(strings.NewReader(mapping)),
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
