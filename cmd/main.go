package main

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/moreirathomas/golastic/internal/repository"
	"github.com/moreirathomas/golastic/pkg/dotenv"
)

const defaultEnvPath = "./.env"

var env = map[string]string{
	"ELASTICSEARCH_INDEX": "",
	"ELASTICSEARCH_SETUP": "",
}

func main() {
	envPath := dotenv.GetPath(defaultEnvPath)

	if err := run(envPath); err != nil {
		log.Fatal(err)
	}
}

func run(envPath string) error {
	if err := dotenv.Load(envPath, &env); err != nil {
		return err
	}

	return initClient()
}

func initClient() error {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return fmt.Errorf("error creating Elasticsearch client: %s", err)
	}

	cfg := repository.Config{
		Client:    client,
		IndexName: env["ELASTICSEARCH_INDEX"],
	}

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		return fmt.Errorf("error creating the repository: %s", err)
	}

	if env["ELASTICSEARCH_SETUP"] == "true" {
		log.Println("Creating Elasticsearch index with mapping")
		if err := setupIndex(repo); err != nil {
			return err
		}
	}

	if err := getESClientInfo(repo); err != nil {
		return err
	}

	return nil
}

func setupIndex(repo *repository.Repository) error {
	// TODO this may be extracted inside a json file and read when needed.
	mapping := `{
	"mappings": {
		"properties": {
			"id":         { "type": "keyword" },
			"title":      { "type": "text", "analyzer": "english" },
			"asbtract":   { "type": "text", "analyzer": "english" }
		}
	}}`

	return repo.CreateIndex(mapping)
}

func getESClientInfo(repo *repository.Repository) error {
	res, err := repo.Info()
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}
