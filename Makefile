GO := go
BINARY := bin/platform-service
ENV_FILE := deployments/docker/.env
COMPOSE := docker compose --env-file $(ENV_FILE) -f deployments/docker/compose.yaml

.PHONY: help fmt test run run-postgres build check clean postgres-up postgres-down postgres-status

help:
	@echo "Available commands:"
	@echo "  make fmt              - Format all Go code"
	@echo "  make test             - Run all tests"
	@echo "  make run              - Run the API with SQLite"
	@echo "  make run-postgres     - Run the API with PostgreSQL"
	@echo "  make build            - Build the API binary"
	@echo "  make check            - Format and test the project"
	@echo "  make clean            - Remove generated build files"
	@echo "  make postgres-up      - Start PostgreSQL"
	@echo "  make postgres-down    - Stop PostgreSQL"
	@echo "  make postgres-status  - Show PostgreSQL status"

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

build:
	mkdir -p bin
	$(GO) build -o $(BINARY) ./cmd/api

check: fmt test

clean:
	rm -rf bin

postgres-up:
	$(COMPOSE) up -d

postgres-down:
	$(COMPOSE) down

postgres-status:
	$(COMPOSE) ps