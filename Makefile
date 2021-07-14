# Run server

.PHONY: run
run:
	@go run ./cmd/main.go

# Docker

.PHONY: docker
docker:
	@docker-compose --env-file ./.env up --build

.PHONY: docker-down
docker-down:
	@docker-compose -env-file ./.env down

.PHONY: docker-flush
docker-flush:
	@docker-compose -env-file ./.env down --volumes

.PHONY: elasticsearch
elasticsearch:
	@docker-compose --env-file ./.env up --build elasticsearch

.PHONY: kibana
kibana:
	@docker-compose --env-file ./.env up --build kibana

# Lint commends

.PHONY: lint
lint:
	@golangci-lint run

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
