package migrations_test

import (
	"testing"

	"github.com/MusdafOmar/platform-service/internal/database"
	"github.com/MusdafOmar/platform-service/migrations"
)

func TestApplyCreatesServicesTable(t *testing.T) {
	db, err := database.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open SQLite database: %v", err)
	}
	defer db.Close()

	if err := migrations.Apply(db); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	var tableName string

	err = db.QueryRow(`
		SELECT name
		FROM sqlite_master
		WHERE type = 'table' AND name = 'services'
	`).Scan(&tableName)
	if err != nil {
		t.Fatalf("find services table: %v", err)
	}

	if tableName != "services" {
		t.Fatalf("expected services table, got %q", tableName)
	}
}
