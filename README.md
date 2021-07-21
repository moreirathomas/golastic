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

When you're done, you can shut down the containers with the following command:

```sh
make down

# equivalent to:
docker-compose down
```

### Local development

#### Prerequisites

- A working version of [Docker](https://docs.docker.com/get-docker/) for elasticsearch and kibana instances

- [Go 1.16](https://golang.org/doc/install) minimum is required due to the use of lastest features, such as `go:embed`

- [golangci-lint](https://golangci-lint.run/) is recommended to run the linters before pushing:

  ```sh
  make lint

  # equivalent to:
  golangci-lint run
  ```

#### Run the local server

The following command runs Elasticsearch containers in the background then starts the server:

```sh
make local

# equivalent to:
docker-compose --env-file ./.env.local up --detach elasticsearch kibana && \
go run cmd/main.go --env-file ./.env.local
```

Alternatively, if you wish to keep track of elasticsearch containers output:

```sh
# Terminal window 1 - run elasticsearch containers
make local-env

# equivalent to:
docker-compose --env-file ./.env.local up elasticsearch kibana


# Terminal window 2 - run web server
make local-server

# equivalent to:
go run cmd/main.go --env-file ./.env.local
```

#### Test & lint

Run all tests:

```sh
make test

# equivalent to:
go test -v -timeout 30s ./... 
```

Run specific test:

```sh
make test t=TestMarshaling p=pkg/golastic

# equivalent to:
go test -v -timeout 30s -run TestMarshaling ./pkg/golastic
```

Run linter (golangci-lint is required):

```sh
make lint

# equivalent to:
golangci-lint run
```

### Populating Elasticsearch

For now, we use the `-p` flag on start-up to populate Elasticsearch:

```sh
go run cmd/main.go -p
```

### Queries

:construction: WIP
