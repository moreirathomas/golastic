# Golastic

Golastic is a web API offering full text search on a books collection with Elasticsearch.

## Get started

### Prerequisites

- For development, Go 1.16 is required
- A working version of Docker with docker-compose is required
- You must provide a `.env` file inside the root directory.
  For a quick start, you can use the values from the provided example:

  ```sh
  echo "$(cat .env.example)" >> .env
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
