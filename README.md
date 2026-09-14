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
- CockroachDB support through the PostgreSQL protocol
- Automatic database migrations
- Docker Compose for PostgreSQL and CockroachDB
- Unit and HTTP handler tests
- Makefile commands for local development

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

Stop the API with:

```text
Control + C
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

Stop the API with `Control + C`, then stop PostgreSQL:

```bash
make postgres-down
```

The named Docker volume preserves the PostgreSQL data when the container is stopped.

## Run with CockroachDB

Start CockroachDB:

```bash
make cockroach-up
```

Create the application database:

```bash
make cockroach-init
```

Check its status:

```bash
make cockroach-status
```

Run the API with CockroachDB:

```bash
make run-cockroach
```

The CockroachDB administration interface is available at:

```text
http://localhost:8081
```

Stop the API with `Control + C`, then stop CockroachDB:

```bash
make cockroach-down
```

The named Docker volume preserves the CockroachDB data when the container is stopped.

## Database configuration

The application selects its database using environment variables:

```text
DB_DRIVER
DB_DSN
```

Default SQLite configuration:

```text
DB_DRIVER=sqlite
DB_DSN=platform-service.db
```

PostgreSQL and CockroachDB both use the pgx driver:

```text
DB_DRIVER=pgx
```

The DSN identifies the database server, port, user, database name and connection options.

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

A successful request returns:

```text
201 Created
```

## List services

```bash
curl http://localhost:8080/services
```

A successful request returns:

```text
200 OK
```

## Run tests

```bash
make test
```

Or format the code and run all tests:

```bash
make check
```

## Available Make commands

```bash
make help
make fmt
make test
make run
make run-postgres
make run-cockroach
make build
make check
make clean
make postgres-up
make postgres-down
make postgres-status
make cockroach-up
make cockroach-init
make cockroach-down
make cockroach-status
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
│   ├── docker/
│   │   ├── .env.example
│   │   ├── compose.yaml
│   │   └── cockroach.compose.yaml
│   └── gocd/
├── docs/
├── infrastructure/
│   └── terraform/
├── internal/
│   ├── config/
│   ├── database/
│   │   ├── database.go
│   │   └── database_test.go
│   ├── server/
│   └── service/
│       ├── repository.go
│       ├── repository_test.go
│       ├── list.go
│       └── list_test.go
├── migrations/
│   ├── 001_create_services.sql
│   ├── migrations.go
│   └── migrations_test.go
├── scripts/
├── tests/
│   └── integration/
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Database compatibility

The repository uses numbered SQL parameters:

```sql
VALUES ($1, $2, $3, $4)
```

This parameter style works with:

- SQLite
- PostgreSQL
- CockroachDB

## Planned milestones

1. Local development environment — complete
2. Google Cloud foundation — complete
3. Go service foundation — complete
4. SQLite database integration — complete
5. PostgreSQL integration — complete
6. Database configuration and CockroachDB compatibility — complete
7. Dockerize the Go API — next
8. Terraform infrastructure
9. Integration testing
10. GoCD pipeline
11. Google Cloud deployment