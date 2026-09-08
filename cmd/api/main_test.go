package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	healthHandler(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.StatusCode,
		)
	}

	contentType := response.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf(
			"expected Content-Type application/json, got %q",
			contentType,
		)
	}

	expectedBody := `{"status":"ok","service":"platform-service"}`
	actualBody := strings.TrimSpace(recorder.Body.String())

	if actualBody != expectedBody {
		t.Errorf(
			"expected body %q, got %q",
			expectedBody,
			actualBody,
		)
	}
}
