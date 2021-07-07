package main

import (
	"log"

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
