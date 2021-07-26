# Docker build

.PHONY: docker
docker:
	@docker-compose --env-file ./.env.docker up --build

.PHONY: down
down:
	@docker-compose down

# Local development

.PHONY: local
local: # Runs Elasticsearch containers in the background then runs the server
	@echo "Starting Elasticsearch containers..." && \
	docker-compose --env-file ./.env.local up --detach elasticsearch kibana && \
	echo "Elasticsearch containers ready." && \
	echo "Starting local server..." && \
	make local-server

.PHONY: local-env
local-env:
	@docker-compose --env-file ./.env.local up elasticsearch kibana

.PHONY: local-server
local-server:
	@go run ./cmd/main.go --env-file ./.env.local

# Clean-up temporary data

.PHONY: clear
clear:
	@rm -rf .logs/* .volumes/*

# Lint commends

.PHONY: lint
lint:
	@golangci-lint run

# Test commands

TEST_FUNC=^.*$$
ifdef t
TEST_FUNC=$(t)
endif
TEST_PKG=./...
ifdef p
TEST_PKG=./$(p)
endif

.PHONY: test
test:
	go test -v -timeout 30s -run $(TEST_FUNC) $(TEST_PKG)

# Serve docs

.PHONY: docs
docs:
	@godoc -http=localhost:9995
