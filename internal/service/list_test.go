package service

import (
	"context"
	"testing"

	"github.com/MusdafOmar/platform-service/internal/database"
	"github.com/MusdafOmar/platform-service/migrations"
)

func TestRepositoryList(t *testing.T) {
	db, err := database.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open SQLite database: %v", err)
	}
	defer db.Close()

	if err := migrations.Apply(db); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	repository := NewRepository(db)
	ctx := context.Background()

	_, err = repository.Create(
		ctx,
		"catalog-api",
		"Manages the service catalog",
	)
	if err != nil {
		t.Fatalf("create first service: %v", err)
	}

	_, err = repository.Create(
		ctx,
		"billing-api",
		"Manages billing",
	)
	if err != nil {
		t.Fatalf("create second service: %v", err)
	}

	services, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("list services: %v", err)
	}

	if len(services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(services))
	}

	if services[0].Name != "catalog-api" {
		t.Fatalf(
			"expected first service catalog-api, got %q",
			services[0].Name,
		)
	}

	if services[1].Name != "billing-api" {
		t.Fatalf(
			"expected second service billing-api, got %q",
			services[1].Name,
		)
	}
}
