# Run

.PHONY: run
run:
	@go run ./cmd/main.go

# Test commands

.PHONY: test
test:
	@go test -v -timeout 30s -run ${t} ./...

.PHONY: tests
tests:
	@go test -v -timeout 30s ./...

# Serve docs

.PHONY: docs
docs:
	@godoc -http=localhost:9995
