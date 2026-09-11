package main

import (
	"encoding/json"
	"net/http"

	"github.com/MusdafOmar/platform-service/internal/service"
)

type listServicesResponse struct {
	Services []service.Service `json:"services"`
}

func listServicesHandler(
	repository *service.Repository,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := repository.List(r.Context())
		if err != nil {
			writeJSONError(
				w,
				"could not list services",
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(
			listServicesResponse{Services: services},
		); err != nil {
			return
		}
	}
}
