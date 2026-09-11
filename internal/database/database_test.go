package database

import "testing"

func TestOpenSQLite(t *testing.T) {
	db, err := Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("expected SQLite connection to open: %v", err)
	}
	defer db.Close()

	var result int

	if err := db.QueryRow("SELECT 1").Scan(&result); err != nil {
		t.Fatalf("expected test query to succeed: %v", err)
	}

	if result != 1 {
		t.Fatalf("expected result 1, got %d", result)
	}
}
