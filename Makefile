GO := go
BINARY := bin/platform-service

ENV_FILE := deployments/docker/.env
POSTGRES_COMPOSE := docker compose --env-file $(ENV_FILE) -f deployments/docker/compose.yaml

COCKROACH_COMPOSE := docker compose -f deployments/docker/cockroach.compose.yaml
COCKROACH_DSN := postgresql://root@localhost:26257/platform_service?sslmode=disable

API_COMPOSE := docker compose -f deployments/docker/api.compose.yaml
API_IMAGE := platform-service:local

.PHONY: help fmt test run run-postgres run-cockroach build check clean \
	docker-build api-up api-down api-status api-logs \
	postgres-up postgres-down postgres-status \
	cockroach-up cockroach-init cockroach-down cockroach-status

help:
	@echo "Available commands:"
	@echo "  make fmt               - Format all Go code"
	@echo "  make test              - Run all tests"
	@echo "  make run               - Run the API with SQLite"
	@echo "  make run-postgres      - Run the API with PostgreSQL"
	@echo "  make run-cockroach     - Run the API with CockroachDB"
	@echo "  make build             - Build the API binary"
	@echo "  make check             - Format and test the project"
	@echo "  make clean             - Remove generated build files"
	@echo "  make docker-build      - Build the API Docker image"
	@echo "  make api-up            - Build and start the containerized API"
	@echo "  make api-down          - Stop the containerized API"
	@echo "  make api-status        - Show the containerized API status"
	@echo "  make api-logs          - Follow the containerized API logs"
	@echo "  make postgres-up       - Start PostgreSQL"
	@echo "  make postgres-down     - Stop PostgreSQL"
	@echo "  make postgres-status   - Show PostgreSQL status"
	@echo "  make cockroach-up      - Start CockroachDB"
	@echo "  make cockroach-init    - Create the CockroachDB database"
	@echo "  make cockroach-down    - Stop CockroachDB"
	@echo "  make cockroach-status  - Show CockroachDB status"

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

run:
	$(GO) run ./cmd/api

run-postgres:
	@set -a; . $(ENV_FILE); set +a; \
	DB_DRIVER=pgx \
	DB_DSN="postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@localhost:$${POSTGRES_PORT}/$${POSTGRES_DB}?sslmode=disable" \
	$(GO) run ./cmd/api

run-cockroach:
	DB_DRIVER=pgx DB_DSN="$(COCKROACH_DSN)" $(GO) run ./cmd/api

build:
	mkdir -p bin
	$(GO) build -o $(BINARY) ./cmd/api

check: fmt test

clean:
	rm -rf bin

docker-build:
	docker build -t $(API_IMAGE) .

api-up:
	$(API_COMPOSE) up -d --build

api-down:
	$(API_COMPOSE) down

api-status:
	$(API_COMPOSE) ps

api-logs:
	$(API_COMPOSE) logs -f api

postgres-up:
	$(POSTGRES_COMPOSE) up -d

postgres-down:
	$(POSTGRES_COMPOSE) down

postgres-status:
	$(POSTGRES_COMPOSE) ps

cockroach-up:
	$(COCKROACH_COMPOSE) up -d

cockroach-init:
	$(COCKROACH_COMPOSE) exec cockroach cockroach sql \
		--insecure \
		--host=localhost:26257 \
		--execute="CREATE DATABASE IF NOT EXISTS platform_service;"

cockroach-down:
	$(COCKROACH_COMPOSE) down

cockroach-status:
	$(COCKROACH_COMPOSE) ps