package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/MusdafOmar/platform-service/internal/database"
	"github.com/MusdafOmar/platform-service/internal/service"
	"github.com/MusdafOmar/platform-service/migrations"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func main() {
	dataSourceName := environmentOrDefault(
		"DB_DSN",
		"platform-service.db",
	)

	db, err := database.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer db.Close()

	if err := migrations.Apply(db); err != nil {
		log.Fatalf("apply database migrations: %v", err)
	}

	serviceRepository := service.NewRepository(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc(
		"POST /services",
		createServiceHandler(serviceRepository),
	)

	address := ":8080"

	log.Printf("database connected: %s", dataSourceName)
	log.Printf("platform-service is running on http://localhost%s", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	response := healthResponse{
		Status:  "ok",
		Service: "platform-service",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}

func environmentOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
