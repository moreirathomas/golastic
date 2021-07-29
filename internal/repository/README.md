# Repository

The package `repository` implements all services related to the data access.

It uses the official Go client for Elasticsearch (`elastic/go-elasticsearch`) and our custom API interface package [`golastic`](../../pkg/golastic/README.md).

The package also provides its own error definitions and methods in `error.go`.
