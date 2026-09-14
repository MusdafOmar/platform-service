# Platform Service

A learning-focused platform and service project built for LIA preparation at Kaerus Software.

The project combines:

- Go REST API
- SQLite, PostgreSQL and CockroachDB
- Docker Compose
- Terraform on Google Cloud
- Makefile automation
- Integration tests
- GoCD continuous delivery

## Current status

The API currently provides:

- `GET /health`
- `POST /services`
- `GET /services`
- SQLite support
- PostgreSQL support through pgx
- Automatic database migrations
- Docker Compose for local PostgreSQL
- Unit and HTTP handler tests

## Requirements

- Go 1.27 or later
- GNU Make
- Git
- Docker Desktop
- Terraform
- Google Cloud CLI

## Run with SQLite

SQLite is the default database:

```bash
make run
```

The local database is stored in:

```text
platform-service.db
```

## Run with PostgreSQL

Create the private environment file:

```bash
cp deployments/docker/.env.example deployments/docker/.env
```

Start PostgreSQL:

```bash
make postgres-up
```

Check its status:

```bash
make postgres-status
```

Run the API with PostgreSQL:

```bash
make run-postgres
```

Stop PostgreSQL:

```bash
make postgres-down
```

The named Docker volume preserves the PostgreSQL data when the container is stopped.

## API address

The API starts at:

```text
http://localhost:8080
```

## Test the health endpoint

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok","service":"platform-service"}
```

## Create a service

```bash
curl \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"catalog-api","description":"LIA platform service"}' \
  http://localhost:8080/services
```

## List services

```bash
curl http://localhost:8080/services
```

## Available Make commands

```bash
make help
make fmt
make test
make run
make run-postgres
make build
make check
make clean
make postgres-up
make postgres-down
make postgres-status
```

## Project structure

```text
platform-service/
├── cmd/
│   └── api/
│       ├── main.go
│       ├── main_test.go
│       ├── services.go
│       ├── services_test.go
│       ├── list_services.go
│       └── list_services_test.go
├── deployments/
│   └── docker/
│       ├── .env.example
│       └── compose.yaml
├── internal/
│   ├── database/
│   └── service/
├── migrations/
│   ├── 001_create_services.sql
│   ├── migrations.go
│   └── migrations_test.go
├── tests/
│   └── integration/
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Planned milestones

1. Local development environment — complete
2. Google Cloud foundation — complete
3. Go service foundation — complete
4. SQLite database integration — complete
5. PostgreSQL integration — complete
6. Database configuration and CockroachDB compatibility — next
7. Dockerize the Go API
8. Terraform infrastructure
9. Integration testing
10. GoCD pipeline
11. Google Cloud deployment