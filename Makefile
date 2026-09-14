GO := go
BINARY := bin/platform-service

ENV_FILE := deployments/docker/.env
POSTGRES_COMPOSE := docker compose --env-file $(ENV_FILE) -f deployments/docker/compose.yaml

COCKROACH_COMPOSE := docker compose -f deployments/docker/cockroach.compose.yaml
COCKROACH_DSN := postgresql://root@localhost:26257/platform_service?sslmode=disable

.PHONY: help fmt test run run-postgres run-cockroach build check clean \
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