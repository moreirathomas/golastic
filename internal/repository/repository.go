package repository

import (
	"errors"
	"fmt"
	"log"

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

func (r Repository) Context() golastic.ContextConfig {
	return golastic.ContextConfig{
		IndexName: r.indexName,
		Client:    r.es,
	}
}

// New returns a new instance of repository.
func New(cfg Config) (*Repository, error) {
	if cfg.IndexName == "" {
		return &Repository{}, errors.New("cannot use empty string \"\" as index name")
	}

	repo := Repository{
		es:        cfg.Client,
		indexName: cfg.IndexName,
	}

	if err := repo.setupIndex(cfg.Mapping); err != nil {
		return nil, err
	}

	return &repo, nil
}

func (r *Repository) setupIndex(mapping string) error {
	isCreate, err := golastic.Indices(r.es).CreateIfNotExists(r.indexName, mapping)
	if isCreate {
		log.Println("Creating Elasticsearch index with mapping")
	}
	if err != nil {
		return fmt.Errorf("cannot create index: %s", err)
	}

	return nil
}

// Info returns basic information about the Elasticsearch client.
func (r *Repository) Info() (*esapi.Response, error) {
	res, err := r.es.Info()
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return res, nil
}
