# Golastic

Golastic is a web API offering full text search on a books collection with Elasticsearch.

## Get started

### Simple usage

- Make sure you have [Docker](https://docs.docker.com/get-docker/) installed.
- Build the app with the following command:
  ```sh
  make

  # equivalent to:
	docker-compose --env-file ./.env.docker up --build
  ```

That's it! The app is ready to be used.
See [Queries](#queries) to learn about our available queries.

### Local development

#### Prerequisites

- A working version of [Docker](https://docs.docker.com/get-docker/) for elasticsearch and kibana instances
- [Go 1.16](https://golang.org/doc/install) minimum is required due to the use of lastest features, such as `go:embed`
- [golangci-lint](https://golangci-lint.run/) is recommended to run the linters before pushing:
  ```sh
  make lint

  # or
  golangci-lint run
  ```

#### Run the local server

- run Elasticsearch and Kibana instances

  ```sh
  make local

  # equivalent to:
  docker-compose --env-file ./.env.local up --build
  ```

- run the server locally

  ```sh
  make local-server

  # equivalent to:
  go run cmd/main.go
  ```

### Populating Elasticsearch

For now, we use flags on start-up to populate Elasticsearch:

```sh
go run cmd/main.go -p
```

### Stop and clean up

```sh
make docker-flush

# or
docker-compose -env-file ./.env down --volumes
```
