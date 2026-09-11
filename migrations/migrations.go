package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
)

//go:embed *.sql
var migrationFiles embed.FS

func Apply(db *sql.DB) error {
	files, err := fs.Glob(migrationFiles, "*.sql")
	if err != nil {
		return fmt.Errorf("find migrations: %w", err)
	}

	sort.Strings(files)

	for _, filename := range files {
		query, err := migrationFiles.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", filename, err)
		}

		if _, err := db.Exec(string(query)); err != nil {
			return fmt.Errorf("apply migration %s: %w", filename, err)
		}
	}

	return nil
}
