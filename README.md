# Golastic

Golastic is a web API offering full text search on a books collection with Elasticsearch.

## Get started

### Prerequisites

- Go 1.16 is required ([for now](https://github.com/moreirathomas/golastic/issues/17)).
- A working version of Docker with docker-compose is required
- You must provide a `.env` file inside the root directory.
    For a quick start, you can use the values from the provided example:

    ```sh
    echo "$(cat .env.example)" >> .env
    ```

#### Development-only prerequisites

- Go 1.16 is required
- [golangci-lint](https://golangci-lint.run/) is recommended to run the linters
  ```sh
  make lint

  # or
  golangci-lint run
  ```

### Run the app

- run Elasticsearch and Kibana instances

    ```sh
    make docker

    # or
    docker-compose --env-file ./.env up --build
    ```

- run the Go server

    ```sh
    make

    # or
    go run cmd/main.go
    ```

### Queries

For now, we use flags on start-up to populate the database and to run queries:

```sh
# populate
go run cmd/main.go -p

# search query
go run cmd/main.go -q "foo"
```

### Stop and clean up

```sh
make docker-flush

# or
docker-compose -env-file ./.env down --volumes
```
