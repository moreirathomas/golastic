package main

import (
	"flag"
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
	query := flag.String("q", "foo", "String value used to search for a match")
	flag.Parse()

	envPath := dotenv.GetPath(defaultEnvPath)

	if err := run(envPath, *query); err != nil {
		log.Fatal(err)
	}
}

func run(envPath string, query string) error {
	if err := dotenv.Load(envPath, &env); err != nil {
		return err
	}

	return initClient(query)
}

func initClient(query string) error {
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

	if err := executeSearch(repo, query); err != nil {
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

func executeSearch(repo *repository.Repository, query string) error {
	res, err := repo.Search(query)
	if err != nil {
		return err
	}

	log.Println(res.Total)
	for _, hit := range res.Hits {
		log.Printf("%#v", hit)
	}

	return nil
}
