# Docker build

.PHONY: docker
docker:
	@docker-compose --env-file ./.env.docker up --build

.PHONY: docker-down
docker-down:
	@docker-compose --env-file ./.env.docker down

# Local develpment

.PHONY: local
local:
	@docker-compose --env-file ./.env.local up --build elasticsearch kibana

.PHONY: local-down
local-down:
	@docker-compose --env-file ./.env.local down

.PHONY: local-server
local-server:
	@go run ./cmd/main.go --env-file ./.env.local

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
