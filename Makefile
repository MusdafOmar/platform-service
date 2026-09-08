GO := go
BINARY := bin/platform-service

.PHONY: help fmt test run build check clean

help:
	@echo "Available commands:"
	@echo "  make fmt    - Format all Go code"
	@echo "  make test   - Run all tests"
	@echo "  make run    - Run the API locally"
	@echo "  make build  - Build the API binary"
	@echo "  make check  - Format and test the project"
	@echo "  make clean  - Remove generated build files"

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

run:
	$(GO) run ./cmd/api

build:
	mkdir -p bin
	$(GO) build -o $(BINARY) ./cmd/api

check: fmt test

clean:
	rm -rf bin