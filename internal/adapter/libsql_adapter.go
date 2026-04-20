package adapter

import (
	"context"
	"database/sql"

	_ "github.com/tursodatabase/go-libsql"
)

type LibSQLAdapter struct {
	db *sql.DB
}

// Bootstrap creates the schema necessary in LibSQL if it doesn't already exist.
func (a *LibSQLAdapter) Bootstrap(ctx context.Context) error {
		query := `
		CREATE TABLE IF NOT EXISTS snippets (
				id TEXT PRIMARY KEY,
				title TEXT NOT NULL,
				owner_id TEXT NOT NULL
				created_at INTEGER NOT NULL
		);`

		_, err := a.db.ExecContext(ctx, query)
		return err
}

func NewLibSQLAdapter(url string) (*LibSQLAdapter, error) {
	db, err := sql.Open("libsql", url)
	if err != nil {
			return nil, err
	}

	adapter := &LibSQLAdapter{db: db}

	// We use context.Background() here because this all happens at startup.
	if err := adapter.Bootstrap(context.Background()); err != nil {
			return nil, err
	}

	return adapter, nil
}