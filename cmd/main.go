package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/moreirathomas/golastic/internal"
	"github.com/moreirathomas/golastic/internal/repository"
	"github.com/moreirathomas/golastic/pkg/dotenv"
)

const defaultEnvPath = "./.env"

var env = map[string]string{
	"ELASTICSEARCH_INDEX": "",
	"ELASTICSEARCH_SETUP": "",
}

// MockupConfig regroups only flags that will be provided on
// the client's http request.
type MockupConfig struct {
	query    string
	populate bool
}

func main() {
	query := flag.String("q", "foo", "String value used to search for a match")
	populate := flag.Bool("p", false, "Populated Elasticsearch with mockup data")
	flag.Parse()

	envPath := dotenv.GetPath(defaultEnvPath)

	cfg := MockupConfig{
		query:    *query,
		populate: *populate,
	}

	if err := run(envPath, cfg); err != nil {
		log.Fatal(err)
	}
}

func run(envPath string, c MockupConfig) error {
	if err := dotenv.Load(envPath, &env); err != nil {
		return err
	}

	return initClient(c)
}

func initClient(c MockupConfig) error {
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

	if c.populate {
		log.Println("Populating Elasticsearch with mockup data")
		if err := populateWithMockup(repo); err != nil {
			return err
		}
	}

	if err := executeSearch(repo, c.query); err != nil {
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

func populateWithMockup(repo *repository.Repository) error {
	books := []internal.Book{
		{Title: "Foo", Abstract: "Lorem ispum foo", ID: "1"},
		{Title: "Bar", Abstract: "Lorem ispum bar", ID: "2"},
		{Title: "Baz", Abstract: "Lorem ispum baz but with foo also", ID: "3"},
	}

	for _, book := range books {
		if err := repo.Create(book); err != nil {
			return err
		}
	}

	return nil
}
