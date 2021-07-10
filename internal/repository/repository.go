package repository

import (
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/moreirathomas/golastic/pkg/golastic"
)

// Config configures the repository.
type Config struct {
	Client    *elasticsearch.Client
	IndexName string
	Mapping   string
}

// Repository allows to index and search documents.
type Repository struct {
	es        *elasticsearch.Client
	indexName string
}

// New returns a new instance of repository.
func New(c Config) (*Repository, error) {
	if c.IndexName == "" {
		return &Repository{}, errors.New("cannot use empty string \"\" as index name")
	}

	repo := Repository{es: c.Client, indexName: c.IndexName}

	cfg := golastic.ContextConfig{
		IndexName: repo.indexName,
		Client:    repo.es,
	}

	if err := golastic.CreateIndexIfNotExists(cfg, c.Mapping); err != nil {
		return &Repository{}, fmt.Errorf("cannot create index: %s", err)
	}

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
