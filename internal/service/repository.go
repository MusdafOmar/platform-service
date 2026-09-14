package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrNameRequired = errors.New("service name is required")

type Service struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (repository *Repository) Create(
	ctx context.Context,
	name string,
	description string,
) (Service, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	if name == "" {
		return Service{}, ErrNameRequired
	}

	record := Service{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now().UTC(),
	}

	_, err := repository.db.ExecContext(
		ctx,
		`
			INSERT INTO services (id, name, description, created_at)
			VALUES ($1, $2, $3, $4)
		`,
		record.ID,
		record.Name,
		record.Description,
		record.CreatedAt,
	)
	if err != nil {
		return Service{}, fmt.Errorf("insert service: %w", err)
	}

	return record, nil
}
