package main

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/moreirathomas/golastic/internal/repository"
)

func main() {
	if err := initClient(); err != nil {
		log.Fatal(err)
	}
}

func initClient() error {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return fmt.Errorf("error creating Elasticsearch client: %s", err)
	}

	cfg := repository.Config{
		Client:    client,
		IndexName: "books",
	}

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		return fmt.Errorf("error creating the repository: %s", err)
	}

	if err := getESClientInfo(repo); err != nil {
		return err
	}

	return nil
}

func getESClientInfo(repo *repository.Repository) error {
	res, err := repo.Info()
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}
