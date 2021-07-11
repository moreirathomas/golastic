package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/moreirathomas/golastic/internal"
	"github.com/moreirathomas/golastic/internal/http"
	"github.com/moreirathomas/golastic/internal/repository"
	"github.com/moreirathomas/golastic/pkg/dotenv"
)

const defaultEnvPath = "./.env"

var env = map[string]string{
	"ELASTICSEARCH_INDEX": "",
	"SERVER_PORT":         "",
}

//go:embed mapping.json
var mapping string

func main() {
	// TODO temporary flags
	populate := flag.Bool("p", false, "Populated Elasticsearch with mockup data")
	flag.Parse()

	envPath := dotenv.GetPath(defaultEnvPath)

	if err := run(envPath, *populate); err != nil {
		log.Fatal(err)
	}
}

func run(envPath string, populate bool) error {
	if err := dotenv.Load(envPath, &env); err != nil {
		return err
	}

	repo, err := initClient()
	if err != nil {
		return err
	}

	if populate {
		log.Println("Populating Elasticsearch with mockup data")
		if err := populateWithMockup(repo); err != nil {
			return err
		}
	}

	addr := ":" + env["SERVER_PORT"]
	srv := http.NewServer(addr, *repo)
	return srv.Start()
}

func initClient() (*repository.Repository, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("error creating Elasticsearch client: %s", err)
	}

	cfg := repository.Config{
		Client:    client,
		IndexName: env["ELASTICSEARCH_INDEX"],
		Mapping:   mapping,
	}

	repo, err := repository.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the repository: %s", err)
	}

	return repo, nil
}

func populateWithMockup(repo *repository.Repository) error {
	books := []internal.Book{
		{Title: "Foo", Abstract: "Lorem ispum foo"},
		{Title: "Bar", Abstract: "Lorem ispum bar"},
		{Title: "Baz", Abstract: "Lorem ispum baz but with foo also"},
	}

	if err := repo.InsertManyBooks(books); err != nil {
		return err
	}

	return nil
}
