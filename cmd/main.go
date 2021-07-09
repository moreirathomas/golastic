package main

import (
	"log"
	"testing"

	"github.com/moreirathomas/golastic/internal"
	"github.com/moreirathomas/golastic/internal/http"
	"github.com/moreirathomas/golastic/internal/repository"
	"github.com/moreirathomas/golastic/pkg/dotenv"
)

const defaultEnvPath = "./.env"

var env = map[string]string{
	"SERVER_PORT": "",
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
	addr := ":" + env["SERVER_PORT"]
	repo := repository.Repository{}
	srv := http.NewServer(addr, repo)
	return srv.Start()
}

func TestValidate(t *testing.T) {
	b := internal.Book{
		Title: "Wesh",
	}
	log.Println(b.Validate())
	if err := b.Validate(); err == nil {
		t.Fatal("Validate should return a non nil error when invalid")
	}
}
