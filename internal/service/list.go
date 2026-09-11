package service

import (
	"context"
	"fmt"
)

func (repository *Repository) List(
	ctx context.Context,
) ([]Service, error) {
	rows, err := repository.db.QueryContext(
		ctx,
		`
			SELECT id, name, description, created_at
			FROM services
			ORDER BY created_at ASC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("query services: %w", err)
	}
	defer rows.Close()

	services := make([]Service, 0)

	for rows.Next() {
		var record Service

		if err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Description,
			&record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan service: %w", err)
		}

		services = append(services, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate services: %w", err)
	}

	return services, nil
}
