GO := go
BINARY := bin/platform-service

ENV_FILE := deployments/docker/.env
POSTGRES_COMPOSE := docker compose --env-file $(ENV_FILE) -f deployments/docker/compose.yaml

COCKROACH_COMPOSE := docker compose -f deployments/docker/cockroach.compose.yaml
COCKROACH_DSN := postgresql://root@localhost:26257/platform_service?sslmode=disable

API_COMPOSE := docker compose -f deployments/docker/api.compose.yaml
API_IMAGE := platform-service:local
TRIVY ?= trivy
TRIVY_SEVERITY ?= HIGH,CRITICAL
GOCD_COMPOSE := docker compose -f deployments/gocd/compose.yaml
TERRAFORM_DIR := infrastructure/terraform
TERRAFORM := terraform -chdir=$(TERRAFORM_DIR)

GCP_PROJECT := lia-platform-musdaf-2026
GCP_REGION := europe-north1
ARTIFACT_REPOSITORY := platform-service
CLOUD_IMAGE_NAME := platform-service
CLOUD_IMAGE_TAG ?= v0.1.0
CLOUD_IMAGE := $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT)/$(ARTIFACT_REPOSITORY)/$(CLOUD_IMAGE_NAME):$(CLOUD_IMAGE_TAG)

.PHONY: help fmt test run run-postgres run-cockroach build check clean \
	docker-build api-up api-down api-status api-logs \
	postgres-up postgres-down postgres-status \
	cockroach-up cockroach-init cockroach-down cockroach-status \
	terraform-fmt terraform-init terraform-validate terraform-plan \
	terraform-apply terraform-output terraform-check \
	artifact-auth cloud-image-push integration-test \
	gocd-up gocd-down gocd-status gocd-logs \
	security-scan security-image-scan

help:
	@echo "Available commands:"
	@echo "  make fmt                - Format all Go code"
	@echo "  make test               - Run all Go tests"
	@echo "  make integration-test   - Run API integration tests"
	@echo "  make run                - Run the API with SQLite"
	@echo "  make run-postgres       - Run the API with PostgreSQL"
	@echo "  make run-cockroach      - Run the API with CockroachDB"
	@echo "  make build              - Build the API binary"
	@echo "  make check              - Format and test the Go project"
	@echo "  make clean              - Remove generated build files"
	@echo "  make docker-build       - Build the local API Docker image"
	@echo "  make api-up             - Build and start the containerized API"
	@echo "  make api-down           - Stop the containerized API"
	@echo "  make api-status         - Show the containerized API status"
	@echo "  make api-logs           - Follow the containerized API logs"
	@echo "  make postgres-up        - Start PostgreSQL"
	@echo "  make postgres-down      - Stop PostgreSQL"
	@echo "  make postgres-status    - Show PostgreSQL status"
	@echo "  make cockroach-up       - Start CockroachDB"
	@echo "  make cockroach-init     - Create the CockroachDB database"
	@echo "  make cockroach-down     - Stop CockroachDB"
	@echo "  make cockroach-status   - Show CockroachDB status"
	@echo "  make terraform-fmt      - Format Terraform configuration"
	@echo "  make terraform-init     - Initialize Terraform and its backend"
	@echo "  make terraform-validate - Validate Terraform configuration"
	@echo "  make terraform-plan     - Preview Terraform infrastructure changes"
	@echo "  make terraform-apply    - Review and apply Terraform changes"
	@echo "  make terraform-output   - Display Terraform outputs"
	@echo "  make terraform-check    - Format and validate Terraform"
	@echo "  make artifact-auth      - Authenticate Docker to Artifact Registry"
	@echo "  make cloud-image-push   - Build and push the amd64 cloud image"
	@echo "  make gocd-up            - Build and start GoCD"
	@echo "  make gocd-down          - Stop GoCD"
	@echo "  make gocd-status        - Show GoCD container status"
	@echo "  make gocd-logs          - Follow GoCD logs"
	@echo "  make security-scan       - Scan source, dependencies, secrets and configuration"
	@echo "  make security-image-scan - Build and scan the local API Docker image"

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

integration-test:
	$(GO) test -count=1 -tags=integration ./tests/integration -v

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
	docker build --pull -t $(API_IMAGE) .

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

terraform-fmt:
	$(TERRAFORM) fmt -recursive

terraform-init:
	$(TERRAFORM) init

terraform-validate:
	$(TERRAFORM) validate

terraform-plan:
	$(TERRAFORM) plan

terraform-apply:
	$(TERRAFORM) apply

terraform-output:
	$(TERRAFORM) output

terraform-check: terraform-fmt terraform-validate

artifact-auth:
	gcloud auth configure-docker $(GCP_REGION)-docker.pkg.dev

cloud-image-push:
	docker buildx build \
		--platform linux/amd64 \
		--tag $(CLOUD_IMAGE) \
		--push \
		.

gocd-up:
	$(GOCD_COMPOSE) up -d --build

gocd-down:
	$(GOCD_COMPOSE) down

gocd-status:
	$(GOCD_COMPOSE) ps

gocd-logs:
	$(GOCD_COMPOSE) logs -f
security-scan:
	$(TRIVY) fs \
		--scanners vuln,misconfig,secret \
		--severity $(TRIVY_SEVERITY) \
		--exit-code 1 \
		.

security-image-scan: docker-build
	$(TRIVY) image \
		--severity $(TRIVY_SEVERITY) \
		--exit-code 1 \
		$(API_IMAGE)