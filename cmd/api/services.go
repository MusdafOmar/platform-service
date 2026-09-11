package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MusdafOmar/platform-service/internal/service"
)

type createServiceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func createServiceHandler(
	repository *service.Repository,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request createServiceRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&request); err != nil {
			writeJSONError(
				w,
				"could not parse request body",
				http.StatusBadRequest,
			)
			return
		}

		record, err := repository.Create(
			r.Context(),
			request.Name,
			request.Description,
		)
		if err != nil {
			if errors.Is(err, service.ErrNameRequired) {
				writeJSONError(
					w,
					err.Error(),
					http.StatusBadRequest,
				)
				return
			}

			writeJSONError(
				w,
				"could not create service",
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_ = json.NewEncoder(w).Encode(record)
	}
}

func writeJSONError(
	w http.ResponseWriter,
	message string,
	status int,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: message,
	})
}
