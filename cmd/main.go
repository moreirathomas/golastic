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

const defaultEnvFile = "./.env.local"

var env = map[string]string{
	"ELASTICSEARCH_INDEX": "",
	"ELASTICSEARCH_URL":   "",
	"SERVER_PORT":         "",
}

// MockupConfig regroups only flags that will be provided on
// the client's http request.
type MockupConfig struct {
	query    string
	populate bool
}

//go:embed mapping.json
var mapping string

func main() {
	envPath := flag.String("env-file", defaultEnvFile, "environment file path")
	query := flag.String("q", "foo", "String value used to search for a match")
	populate := flag.Bool("p", false, "Populated Elasticsearch with mockup data")
	flag.Parse()

	if err := dotenv.Load(*envPath, env); err != nil {
		log.Fatal(err)
	}

	cfg := MockupConfig{
		query:    *query,
		populate: *populate,
	}

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(c MockupConfig) error {
	repo, err := initClient(c)
	if err != nil {
		return err
	}

	addr := ":" + env["SERVER_PORT"]
	srv := http.NewServer(addr, *repo)
	return srv.Start()
}

func initClient(c MockupConfig) (*repository.Repository, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{env["ELASTICSEARCH_URL"]},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Elasticsearch client: %s", err)
	}

	cfg := repository.Config{
		Client:    client,
		IndexName: env["ELASTICSEARCH_INDEX"],
	}

	repo, err := repository.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the repository: %s", err)
	}

	if err := repo.CreateIndexIfNotExists(mapping); err != nil {
		return nil, err
	}

	if c.populate {
		log.Println("Populating Elasticsearch with mockup data")
		if err := populateWithMockup(repo); err != nil {
			return nil, err
		}
	}

	if err := executeSearch(repo, c.query); err != nil {
		return nil, err
	}

	return repo, nil
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
		{Title: "Foo", Abstract: "Lorem ispum foo"},
		{Title: "Bar", Abstract: "Lorem ispum bar"},
		{Title: "Baz", Abstract: "Lorem ispum baz but with foo also"},
	}

	return repo.CreateBulk(books)
}
