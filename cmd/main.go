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

// MockupConfig regroups only flags that will be provided on
// the client's http request.
type MockupConfig struct {
	query    string
	populate bool
}

//go:embed mapping.json
var mapping string

func main() {
	// TODO temporary flags
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

	repo, err := initClient(c)
	if err != nil {
		return err
	}

	addr := ":" + env["SERVER_PORT"]
	srv := http.NewServer(addr, *repo)
	return srv.Start()
}

func initClient(c MockupConfig) (*repository.Repository, error) {
	client, err := elasticsearch.NewDefaultClient()
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

	// if err := printESClientInfo(repo); err != nil {
	// 	return nil, err
	// }

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

func printESClientInfo(repo *repository.Repository) error {
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
		{Title: "Foo", Abstract: "Lorem ispum foo"},
		{Title: "Bar", Abstract: "Lorem ispum bar"},
		{Title: "Baz", Abstract: "Lorem ispum baz but with foo also"},
	}

	if err := repo.CreateBulk(books); err != nil {
		return err
	}

	return nil
}
