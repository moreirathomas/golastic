# Golastic

Golastic is a web API offering full text search and CRUD operations on a book collection via Elasticsearch.

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
docker-compose --env-file ./.env.local up --detach elasticsearch kibana && \
go run cmd/main.go --env-file ./.env.local
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

Run a specific test:

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

We use a CLI flag to conditionally populate Elasticsearch's index on start-up. To use it, simply run:

```sh
go run cmd/main.go -p --env-file ./.env.local
```

Only the first run (or any run following an erasure of the Docker volume) requires the use of this flag, as the dummy data will not be overwritten.

### Test with queries

#### Search books by full text query

Request:

```sh
curl http://localhost:9999/books?query=foo&page=1&size=10
```

Response:

```json
200 OK

{
  "links": {},
  "page": 1,
  "per_page": 10,
  "results": [
    {
      "abstract": "Lorem ispum baz but with foo also",
      "title": "Baz"
    },
    {
      "abstract": "Lorem ispum foo",
      "title": "Foo"
    }
  ],
  "total": 2
}
```

#### Get a book by ID

Request:

```sh
curl  http://localhost:9999/books/nWJ45HoBEwNIQ_UGmi_R
```

Response:

```json
200 OK

{
  "abstract": "Lorem ispum foo",
  "author": {
    "firstname": "John",
    "lastname": "Doe"
  },
  "created_at": "0001-01-01T00:00:00Z",
  "id": "nWJ45HoBEwNIQ_UGmi_R",
  "title": "Foo"
}
```

#### Create a book

Request:

```sh
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"title": "Post", "abstract": "Created recently", "author": {"firstname": "User", "lastname": "User"}}' \
  http://localhost:9999/books
```

Response:

```txt
201 Created
```

#### Update a book

Request:

```sh
curl -X PUT \
  -H "Content-Type: application/json" \
  -d '{"abstract": "Redacted"}' \
  http://localhost:9999/books/nWJ45HoBEwNIQ_UGmi_R
```

Response:

```txt
204 No Content
```

#### Delete a book

Request:

```sh
curl -X DELETE http://localhost:9999/books/nWJ45HoBEwNIQ_UGmi_R
```

Response:

```txt
204 No Content
```
