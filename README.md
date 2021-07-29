# Golastic <!-- omit in toc -->

Golastic is a web API offering full text search and CRUD operations on a book collection via Elasticsearch.

## Table of contents <!-- omit in toc -->

- [Getting started](#getting-started)
  - [Installation and dependencies](#installation-and-dependencies)
  - [Run the project in Docker](#run-the-project-in-docker)
  - [Run the server locally](#run-the-server-locally)
  - [Test and lint](#test-and-lint)
  - [Populate with dummies](#populate-with-dummies)
  - [Test routes with CURL commands](#test-routes-with-curl-commands)
- [Architecture](#architecture)
  - [Folder structure](#folder-structure)
  - [`internal`](#internal)
  - [`pkg`](#pkg)
  - [`cmd`](#cmd)
- [Further documentation](#further-documentation)
- [Acknowledgement and contributors](#acknowledgement-and-contributors)

## Getting started

### Installation and dependencies

You must be able to run Docker and Docker Compose on your local machine to run this app. Refer to the [Get Docker](https://docs.docker.com/get-docker/) and [Install Docker Compose](https://docs.docker.com/compose/install/) docs for their installation.

### Run the project in Docker

The app is fully dockerized. To start all required containers, simply run:

```sh
make
# alias to:
# docker-compose --env-file ./.env.docker up --build
```

That's it! The app is ready to be used.
See [Test with queries](#test-with-queries) to learn about our available queries.

To shut down the containers, run:

```sh
make down
# equivalent to:
# docker-compose down
```

### Run the server locally

[Go 1.16](https://golang.org/doc/install) minimum is required due to the use of newer features (notably `go:embed`).

To start the server locally with Elasticsearch containers in the background, run:

```sh
make local

# alias to:
# docker-compose --env-file ./.env.local up --detach elasticsearch kibana && \
# go run cmd/main.go --env-file ./.env.local
```

Alternatively, if you wish to keep track of elasticsearch containers output:

```sh
# Terminal window 1 - run elasticsearch containers
make local-env
# alias to:
# docker-compose --env-file ./.env.local up elasticsearch kibana


# Terminal window 2 - run web server
make local-server
# alias to:
# go run cmd/main.go --env-file ./.env.local
```

You may also want to populate Elasticsearch index with dummy data. See [Populate with dummies](#populate-with-dummies) for more on this.

### Test and lint

Run all tests:

```sh
make test
# alias to:
# go test -v -timeout 30s ./...
```

Run a specific test with `t` to specify a test and `p` to specify a package (parameters are independent):

```sh
make test t=TestMarshaling p=pkg/golastic
# alias to:
# go test -v -timeout 30s -run TestMarshaling ./pkg/golastic
```

Run the linter:

> We use [golangci-lint](https://golangci-lint.run/) in our CI. It runs on each push to a branch with an open PR.

```sh
make lint
# alias to:
# golangci-lint run
```

### Populate with dummies

You may use a CLI flag to populate Elasticsearch's index on start-up. To use it, simply run:

```sh
go run cmd/main.go -p
```

Only the first run (or any run following an erasure of the Docker volume) requires the use of this flag, as the dummy data will not be overwritten.

### Test routes with CURL commands

Refer to the [routes specifition](internal/http/README.md) for detailed requests queries and responses data. It comes with handy CURL commands to quickly test the routes at runtime.

## Architecture

### Folder structure

The main functional packages for the project are:

```txt
.
├── cmd
├── internal
│   ├── http
│   └── repository
└── pkg
    ├── golastic
    └── ...
```

### `internal`

The main application code. It defines domain related entities at its root and its child packages are grouped by dependencies (data access and transport).

### `pkg`

Purposeful and reusable library code. It does not import any types from `internal` and does not rely on it to work. There is no domain logic inside this directory.

### `cmd`

Main application for the project. A `main` function imports and call code from `internal` or `pkg`. It ties all dependencies together and injects runtime variables.

## Further documentation

Sub-packages have their own documentation when it is relevant. You may refer to these docs:

- [http](internal/http/README.md)
- [repository](internal/repository/README.md)
- [golastic](pkg/golastic/README.md)

Moreover, all types, functions and methods are documented.

## Acknowledgement and contributors

This project is part of an school assignement. Our team members are:

- Gregory Albouy
- Thomas Moreira
- Damien Mathieu
