# Platform Service

A learning-focused platform and service project built for LIA preparation at Kaerus Software.

The project combines:

- Go REST API
- SQLite, PostgreSQL and CockroachDB
- Docker and Docker Compose
- Terraform on Google Cloud
- Makefile automation
- Integration testing
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
- Unit and HTTP handler tests

The platform currently provides:

- Multi-stage Docker image for the Go API
- Docker Compose for the API, PostgreSQL and CockroachDB
- Persistent local Docker volumes
- Docker container health checks
- Local GoCD server and custom GoCD build agent
- GoCD agent with Go, Git, Make and Docker CLI
- Parallel unit and integration-test jobs
- Docker image build stage after successful verification
- Automatic pipeline scheduling for changes on `main`
- Makefile automation for Go, Docker, Terraform and GoCD
- Terraform-managed Google Cloud infrastructure
- Protected GCS remote Terraform state
- Artifact Registry Docker repository
- Versioned `linux/amd64` cloud container image
- Dedicated Cloud Run service account
- Private Cloud Run deployment
- Scale-to-zero with a maximum of one instance
- 100 SEK monthly GCP budget with alert thresholds

## Requirements

- Go 1.27 or later
- GNU Make
- Git
- Docker Desktop
- Terraform
- Google Cloud CLI

## Run the API in Docker

Build the API Docker image:

```bash
make docker-build
```

The image is created with this name:

```text
platform-service:local
```

Build and start the containerized API:

```bash
make api-up
```

Check the container status:

```bash
make api-status
```

The status should eventually show:

```text
healthy
```

Test the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok","service":"platform-service"}
```

Follow the API container logs:

```bash
make api-logs
```

Press `Control + C` to stop following the logs. This does not stop the container.

Stop and remove the API container and its network:

```bash
make api-down
```

The named volume below preserves the SQLite database when the container is removed:

```text
platform-service-api_sqlite_data
```

Do not add `-v` to the Compose down command unless you intentionally want to delete the stored SQLite data.

## Docker image design

The `Dockerfile` uses a multi-stage build.

The builder stage:

1. Uses the Go Alpine image.
2. Downloads the Go dependencies.
3. Compiles the API into a Linux binary.
4. Produces the `platform-service` executable.

The runtime stage:

1. Uses a smaller Alpine Linux image.
2. Installs CA certificates.
3. Creates a non-root application user.
4. Copies only the compiled binary.
5. Stores SQLite data under `/app/data`.
6. Exposes port `8080`.

The application runs as the non-root user:

```text
app
```

This reduces the privileges available inside the running container.

## Docker build context

The `.dockerignore` file prevents unnecessary or sensitive files from being sent to Docker during a build.

Ignored files include:

- Git history
- Editor settings
- Environment files
- Local database files
- Build output
- Terraform state
- Documentation and deployment files that are not required for compilation

This keeps the build context smaller and prevents local secrets from entering the image.

## Run with SQLite without Docker

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

The named Docker volume preserves PostgreSQL data when the container is stopped or removed.

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

The named Docker volume preserves CockroachDB data when the container is stopped or removed.

The local CockroachDB container uses insecure mode for development only. Production environments must use authentication and encrypted connections.

## Database configuration

The application selects its database using environment variables:

```text
DB_DRIVER
DB_DSN
```

### SQLite configuration

Default local configuration:

```text
DB_DRIVER=sqlite
DB_DSN=platform-service.db
```

The Docker container uses:

```text
DB_DRIVER=sqlite
DB_DSN=/app/data/platform-service.db
```

### PostgreSQL configuration

PostgreSQL uses:

```text
DB_DRIVER=pgx
```

The PostgreSQL DSN contains:

- Database username
- Database password
- Database host
- Database port
- Database name
- SSL configuration

### CockroachDB configuration

CockroachDB also uses:

```text
DB_DRIVER=pgx
```

The local development DSN is:

```text
postgresql://root@localhost:26257/platform_service?sslmode=disable
```

CockroachDB can reuse pgx because it supports the PostgreSQL wire protocol.

## Database migrations

Database migrations are stored in:

```text
migrations/
```

The first migration creates the `services` table:

```text
migrations/001_create_services.sql
```

The Go application embeds the SQL migration files into the compiled binary.

When the API starts, it automatically applies the migrations before starting the HTTP server.

The same migration currently works with:

- SQLite
- PostgreSQL
- CockroachDB

## API address

The API starts at:

```text
http://localhost:8080
```

## Health endpoint

Request:

```bash
curl -i http://localhost:8080/health
```

Expected status:

```text
HTTP/1.1 200 OK
```

Expected response:

```json
{"status":"ok","service":"platform-service"}
```

## Create a service

Request:

```bash
curl -i \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"catalog-api","description":"LIA platform service"}' \
  http://localhost:8080/services
```

Expected status:

```text
HTTP/1.1 201 Created
```

Example response:

```json
{
  "id": "generated-uuid",
  "name": "catalog-api",
  "description": "LIA platform service",
  "created_at": "generated-timestamp"
}
```

The API:

1. Validates the request.
2. Generates a UUID.
3. Creates a UTC timestamp.
4. Inserts the service into the selected database.
5. Returns the created record as JSON.

## List services

Request:

```bash
curl -i http://localhost:8080/services
```

Expected status:

```text
HTTP/1.1 200 OK
```

Example response:

```json
{
  "services": [
    {
      "id": "generated-uuid",
      "name": "catalog-api",
      "description": "LIA platform service",
      "created_at": "generated-timestamp"
    }
  ]
}
```

## Run tests

Run the unit and handler tests:

```bash
make test
```

Format the Go code and run the unit and handler tests:

```bash
make check
```

Run the black-box API integration test:

```bash
make integration-test
```

The integration test performs a fresh run every time. It:

- Builds the real API binary
- Creates a temporary SQLite database
- Selects an available local port
- Starts the API as a separate process
- Waits for `GET /health` to return successfully
- Creates a service with `POST /services`
- Retrieves the service with `GET /services`
- Stops the API and removes temporary test files automatically

The complete test coverage includes:

- Health endpoint
- Create-service handler
- List-services handler
- SQLite database connection
- Service repository
- Service listing
- Database migrations
- Full API service lifecycle

## Run GoCD locally

Build and start the GoCD server and custom build agent:

```bash
make gocd-up
```

Open the GoCD interface:

```text
http://localhost:8153/go
```

Check the container status:

```bash
make gocd-status
```

Follow the GoCD logs:

```bash
make gocd-logs
```

The pipeline contains two stages:

```text
verify
├── unit-tests
└── integration-tests

build
└── docker-build
```

The verification jobs run independently. When both pass, GoCD automatically starts the build stage and creates the local image:

```text
platform-service:local
```

Changes pushed to the `main` branch automatically schedule the pipeline.

Stop GoCD:

```bash
make gocd-down
```

The named Docker volumes preserve the GoCD server configuration and agent data.

## Available Make commands

```bash
make help
make fmt
make test
make integration-test
make run
make run-postgres
make run-cockroach
make build
make check
make clean
make docker-build
make api-up
make api-down
make api-status
make api-logs
make postgres-up
make postgres-down
make postgres-status
make cockroach-up
make cockroach-init
make cockroach-down
make cockroach-status
make gocd-up
make gocd-down
make gocd-status
make gocd-logs
```

### Application commands

| Command | Purpose |
|---|---|
| `make run` | Run the API locally with SQLite |
| `make run-postgres` | Run the API locally with PostgreSQL |
| `make run-cockroach` | Run the API locally with CockroachDB |
| `make build` | Compile the Go API |
| `make test` | Run unit and handler tests |
| `make integration-test` | Build and test the complete API lifecycle |
| `make check` | Format the code and run all tests |
| `make clean` | Remove generated build files |

### API Docker commands

| Command | Purpose |
|---|---|
| `make docker-build` | Build the API Docker image |
| `make api-up` | Build and start the containerized API |
| `make api-down` | Stop and remove the API container |
| `make api-status` | Show the API container status |
| `make api-logs` | Follow the API container logs |

### PostgreSQL commands

| Command | Purpose |
|---|---|
| `make postgres-up` | Start PostgreSQL |
| `make postgres-down` | Stop PostgreSQL |
| `make postgres-status` | Show PostgreSQL status |

### CockroachDB commands

| Command | Purpose |
|---|---|
| `make cockroach-up` | Start CockroachDB |
| `make cockroach-init` | Create the application database |
| `make cockroach-down` | Stop CockroachDB |
| `make cockroach-status` | Show CockroachDB status |

### GoCD commands

| Command | Purpose |
|---|---|
| `make gocd-up` | Build and start the GoCD server and agent |
| `make gocd-down` | Stop and remove the GoCD containers |
| `make gocd-status` | Show the GoCD container status |
| `make gocd-logs` | Follow the GoCD logs |

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
│   │   ├── api.compose.yaml
│   │   ├── compose.yaml
│   │   └── cockroach.compose.yaml
│   └── gocd/
│       ├── agent.Dockerfile
│       └── compose.yaml
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
│       └── api_test.go
├── .dockerignore
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md

## Database compatibility

The repository uses numbered SQL parameters:

```sql
VALUES ($1, $2, $3, $4)
```

This parameter style works with:

- SQLite
- PostgreSQL
- CockroachDB

This allows the same repository implementation to support all three databases.

## Docker persistence

There are three database-storage scenarios in the project:

| Environment | Storage |
|---|---|
| Local SQLite | `platform-service.db` |
| Dockerized SQLite | `platform-service-api_sqlite_data` volume |
| PostgreSQL | `docker_postgres_data` volume |
| CockroachDB | `docker_cockroach_data` volume |

A container can be removed and recreated while its named volume remains available.

Running Compose down normally preserves the volume:

```bash
make api-down
```

Deleting the volume is destructive and removes its stored database data.

## Google Cloud foundation

The Google Cloud project is:

```text
lia-platform-musdaf-2026
```

The following foundation work is complete:

- Google Cloud CLI authentication
- Application Default Credentials
- Billing enabled
- Monthly budget of 100 SEK configured
- Budget thresholds at 25%, 50%, 90% and 100%
- Required Google Cloud APIs enabled
- Regional resources configured in `europe-north1`

## Terraform infrastructure

Terraform manages the Google Cloud infrastructure under:

```text
infrastructure/terraform/
```

The configuration currently manages:

- Required Google Cloud APIs
- Artifact Registry Docker repository
- Dedicated Cloud Run service account
- Private Cloud Run service
- Protected GCS bucket for remote Terraform state
- Terraform outputs for repository and service information

Initialize Terraform:

```bash
make terraform-init
```

Format and validate the configuration:

```bash
make terraform-check
```

Preview infrastructure changes:

```bash
make terraform-plan
```

Apply reviewed infrastructure changes:

```bash
make terraform-apply
```

Display Terraform outputs:

```bash
make terraform-output
```

Terraform state is stored remotely in:

```text
gs://lia-platform-musdaf-2026-terraform-state/platform-service/development/
```

The state bucket uses:

- Object versioning
- Uniform bucket-level access
- Public access prevention
- Terraform destroy protection

Private values belong in:

```text
infrastructure/terraform/terraform.tfvars
```

This file is ignored by Git. Use the tracked example as a template:

```text
infrastructure/terraform/terraform.tfvars.example
```

## Artifact Registry

The Docker repository is:

```text
europe-north1-docker.pkg.dev/lia-platform-musdaf-2026/platform-service
```

Authenticate Docker:

```bash
make artifact-auth
```

Build a Cloud Run-compatible `linux/amd64` image and push it:

```bash
make cloud-image-push
```

The default image is:

```text
europe-north1-docker.pkg.dev/lia-platform-musdaf-2026/platform-service/platform-service:v0.1.0
```

A different version can be pushed with:

```bash
make cloud-image-push CLOUD_IMAGE_TAG=v0.1.1
```

## Cloud Run deployment

The Terraform-managed Cloud Run service is:

```text
platform-service
```

It uses:

- A dedicated service account
- Authenticated access
- Scale-to-zero
- Maximum one running instance
- One vCPU
- 512 MiB memory
- Container port 8080
- Image version `v0.1.0`

Display the service URL:

```bash
terraform -chdir=infrastructure/terraform \
  output -raw cloud_run_service_url
```

Store the URL temporarily:

```bash
LIA_SERVICE_URL="$(
  terraform -chdir=infrastructure/terraform \
  output -raw cloud_run_service_url
)"
```

Test the authenticated health endpoint:

```bash
curl \
  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  "$LIA_SERVICE_URL/health"
```

Expected response:

```json
{"status":"ok","service":"platform-service"}
```

Anonymous requests return `403 Forbidden` because the service is private.

## Cloud database limitation

The initial Cloud Run deployment uses SQLite at:

```text
/tmp/platform-service.db
```

Cloud Run container filesystems are ephemeral. Records can disappear when the instance restarts or scales down.

SQLite in Cloud Run is therefore used only to verify that the deployed container works. A durable cloud deployment requires an external database such as PostgreSQL.

## Planned milestones

1. Local development environment — complete
2. Google Cloud foundation — complete
3. Go service foundation — complete
4. SQLite database integration — complete
5. PostgreSQL integration — complete
6. Database configuration and CockroachDB compatibility — complete
7. Dockerize the Go API — complete
8. Terraform infrastructure — complete
9. Integration testing — complete
10. 10. GoCD pipeline — in progress
11. Google Cloud deployment — initial Cloud Run deployment complete