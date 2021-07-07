package repository

import (
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// Config configures the repository.
type Config struct {
	Client    *elasticsearch.Client
	IndexName string
}

// Repository allows to index and search documents.
type Repository struct {
	es        *elasticsearch.Client
	indexName string
}

// NewRepository returns a new instance of repository.
func NewRepository(c Config) (*Repository, error) {
	if c.IndexName == "" {
		return &Repository{}, errors.New("cannot use empty string \"\" as index name")
	}

	repo := Repository{es: c.Client, indexName: c.IndexName}
	return &repo, nil
}

// Info returns basic information about the Elasticsearch client.
func (r *Repository) Info() (*esapi.Response, error) {
	res, err := r.es.Info()
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return res, nil
}
