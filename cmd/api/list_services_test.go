package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MusdafOmar/platform-service/internal/database"
	"github.com/MusdafOmar/platform-service/internal/service"
	"github.com/MusdafOmar/platform-service/migrations"
)

func TestListServicesHandler(t *testing.T) {
	db, err := database.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open SQLite database: %v", err)
	}
	defer db.Close()

	if err := migrations.Apply(db); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	repository := service.NewRepository(db)

	_, err = repository.Create(
		context.Background(),
		"catalog-api",
		"LIA platform service",
	)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}

	handler := listServicesHandler(repository)
	request := httptest.NewRequest(
		http.MethodGet,
		"/services",
		nil,
	)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	var response listServicesResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(response.Services) != 1 {
		t.Fatalf(
			"expected one service, got %d",
			len(response.Services),
		)
	}

	if response.Services[0].Name != "catalog-api" {
		t.Fatalf(
			"expected catalog-api, got %q",
			response.Services[0].Name,
		)
	}
}
