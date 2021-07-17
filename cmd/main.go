package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/estransport"

	"github.com/moreirathomas/golastic/internal"
	"github.com/moreirathomas/golastic/internal/http"
	"github.com/moreirathomas/golastic/internal/repository"
	"github.com/moreirathomas/golastic/pkg/dotenv"
	"github.com/moreirathomas/golastic/pkg/logger"
)

const (
	defaultEnvFile = "./.env.local"
	logPath        = "./.logs/"
)

var env = map[string]string{
	"ELASTICSEARCH_INDEX": "",
	"ELASTICSEARCH_URL":   "",
	"SERVER_PORT":         "",
}

//go:embed mapping.json
var mapping string

func main() {
	envPath := flag.String("env-file", defaultEnvFile, "environment file path")
	populate := flag.Bool("p", false, "Populated Elasticsearch with mockup data")
	flag.Parse()

	if err := dotenv.Load(*envPath, env); err != nil {
		log.Fatal(err)
	}

	if err := run(*populate); err != nil {
		log.Fatal(err)
	}
}

func run(populate bool) error {
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
	srv.ErrorLog = logger.DefaultFile(logPath + "server.errorlog")
	return srv.Start()
}

func initClient() (*repository.Repository, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{env["ELASTICSEARCH_URL"]},
		Logger: &estransport.TextLogger{
			Output: logger.DefaultFile(logPath + "elasticsearch.log").Writer(),
			// EnableRequestBody:  true,
			// EnableResponseBody: true,
		},
	})
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

	return repo.InsertManyBooks(books)
}
