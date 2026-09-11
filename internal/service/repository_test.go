package service

import (
	"context"
	"testing"

	"github.com/MusdafOmar/platform-service/internal/database"
	"github.com/MusdafOmar/platform-service/migrations"
)

func TestRepositoryCreate(t *testing.T) {
	db, err := database.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open SQLite database: %v", err)
	}
	defer db.Close()

	if err := migrations.Apply(db); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	repository := NewRepository(db)

	record, err := repository.Create(
		context.Background(),
		"billing-service",
		"Handles billing operations",
	)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}

	if record.ID == "" {
		t.Fatal("expected generated service ID")
	}

	if record.Name != "billing-service" {
		t.Fatalf(
			"expected name billing-service, got %q",
			record.Name,
		)
	}

	var storedName string
	var storedDescription string

	err = db.QueryRow(
		`
			SELECT name, description
			FROM services
			WHERE id = ?
		`,
		record.ID,
	).Scan(&storedName, &storedDescription)
	if err != nil {
		t.Fatalf("read stored service: %v", err)
	}

	if storedName != record.Name {
		t.Fatalf(
			"expected stored name %q, got %q",
			record.Name,
			storedName,
		)
	}

	if storedDescription != record.Description {
		t.Fatalf(
			"expected stored description %q, got %q",
			record.Description,
			storedDescription,
		)
	}
}
