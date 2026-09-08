# Platform Service

A learning-focused platform and service project built for LIA preparation.

The project will combine:

- Go REST API
- PostgreSQL, SQLite and CockroachDB
- Docker
- Terraform on Google Cloud
- Makefile automation
- Integration tests
- GoCD continuous delivery

## Current status

The service currently provides a health endpoint and an automated test.

## Requirements

- Go 1.27 or later
- GNU Make
- Git
- Docker Desktop
- Terraform
- Google Cloud CLI

## Run locally

```bash
make run
```

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

## Available Make commands

```bash
make help
make fmt
make test
make run
make build
make check
make clean
```

## Project structure

```text
platform-service/
├── cmd/
│   └── api/
│       ├── main.go
│       └── main_test.go
├── .gitignore
├── go.mod
├── Makefile
└── README.md
```

## Planned milestones

1. Local development environment
2. Google Cloud foundation
3. Go service foundation
4. Database integration
5. Docker development environment
6. Terraform infrastructure
7. Integration testing
8. GoCD pipeline
9. Google Cloud deployment