package repository

import (
	"errors"
	"fmt"
	"log"
	"strings"

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

// New returns a new instance of repository.
func New(c Config) (*Repository, error) {
	if c.IndexName == "" {
		return &Repository{}, errors.New("cannot use empty string \"\" as index name")
	}

	repo := Repository{es: c.Client, indexName: c.IndexName}
	return &repo, nil
}

// CreateIndexIfNotExists creates a new index with mapping
// if the repository's unique index does not exists yet.
func (r *Repository) CreateIndexIfNotExists(mapping string) error {
	exists, err := r.indexExists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	log.Println("Creating Elasticsearch index with mapping")
	return r.createIndex(mapping)
}

// indexExists returns true when the index already exists in the repository.
func (r *Repository) indexExists() (bool, error) {
	res, err := r.es.Indices.Exists([]string{r.indexName})
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

// createIndex creates a new index with mapping.
func (r *Repository) createIndex(mapping string) error {
	res, err := r.es.Indices.Create(r.indexName, r.es.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
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
