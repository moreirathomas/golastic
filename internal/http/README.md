# HTTP

The package `http` implements all services related to HTTP transport and routing.

`request.go` and `response.go` provide reusable functions for handling requests and manipulating response objects.

The package also provides its own error definitions and methods in `error.go`.

## Routes specification

> Note: for quick testing, curl commands are provided.

### Search books by full text query

Request:

```sh
curl http://localhost:9999/books?query=<query_string>&page=1&size=10
```

Response:

```json
200 OK

{
   "links" : {
      "next" : "http://localhost:9999/books?query=foo&page=2&size=5"
   },
   "page" : 1,
   "per_page" : 5,
   "results" : [
      {
         "abstract" : "Lorem ispum foo",
         "author" : {
            "firstname" : "John",
            "lastname" : "Doe"
         },
         "created_at" : "2021-07-26T22:34:21.516269+02:00",
         "id" : "oGKG5HoBEwNIQ_UGji_k",
         "title" : "Foo"
      },
      {
         "abstract" : "Lorem ispum bar and foo",
         "author" : {
            "firstname" : "John",
            "lastname" : "Doe"
         },
         "created_at" : "2021-07-27T11:36:03.230521+02:00",
         "id" : "omJS53oBEwNIQ_UGOC-q",
         "title" : "Bar"
      },
      // ...
   ],
   "total" : 10
}
```

### Get a book by ID

Request:

```sh
curl  http://localhost:9999/books/<id>
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
  "created_at": "2021-07-26T22:34:21.516269+02:00",
  "id": "oGKG5HoBEwNIQ_UGji_k",
  "title": "Foo"
}
```

### Create a book

Request:

```sh
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"title": "Create a Book", "abstract": "It is new!", "author": {"firstname": "John", "lastname": "Doe"}}' \
  http://localhost:9999/books
```

Response:

```json
201 Created
{
  "abstract": "It is new!",
  "author": {
    "firstname": "John",
    "lastname": "Doe"
  },
  "created_at": "2021-07-26T22:34:21.516269+02:00",
  "id": "nWJ45HoBEwNIQ_UGmi_R",
  "title": "Create a Book"
}
```

### Update a book

Request:

```sh
curl -X PUT \
  -H "Content-Type: application/json" \
  -d '{"abstract": "It is updated!"}' \
  http://localhost:9999/books/<id>
```

Response:

```txt
204 No Content
```

### Delete a book

Request:

```sh
curl -X DELETE http://localhost:9999/books/<id>
```

Response:

```txt
204 No Content
```
