package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MusdafOmar/platform-service/internal/database"
	"github.com/MusdafOmar/platform-service/internal/service"
	"github.com/MusdafOmar/platform-service/migrations"
)

func TestCreateServiceHandler(t *testing.T) {
	db, err := database.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open SQLite database: %v", err)
	}
	defer db.Close()

	if err := migrations.Apply(db); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	repository := service.NewRepository(db)
	handler := createServiceHandler(repository)

	body := `{
		"name": "catalog-api",
		"description": "LIA platform service"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/services",
		strings.NewReader(body),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusCreated,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	var response service.Service

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.ID == "" {
		t.Fatal("expected generated service ID")
	}

	if response.Name != "catalog-api" {
		t.Fatalf(
			"expected name catalog-api, got %q",
			response.Name,
		)
	}

	var storedCount int

	if err := db.QueryRow(
		"SELECT COUNT(*) FROM services WHERE id = ?",
		response.ID,
	).Scan(&storedCount); err != nil {
		t.Fatalf("count stored service: %v", err)
	}

	if storedCount != 1 {
		t.Fatalf("expected one stored service, got %d", storedCount)
	}
}
