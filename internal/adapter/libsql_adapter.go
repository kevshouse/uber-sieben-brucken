package adapter

import (
	"context"
	"database/sql"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	_ "github.com/ncruces/go-sqlite3/driver" // '_' forces the compiler to keep package for registering 'libsql' at runtime
)

type LibSQLAdapter struct {
	db *sql.DB
}

// NewLibSQLAdapter initializes the connection and prepares the 'Shore'
func NewLibSQLAdapter(url string) (*LibSQLAdapter, error) {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		return nil, err
	}

	adapter := &LibSQLAdapter{db: db}

	// Bootstrap ensures the snippets table exists at startup
	if err := adapter.Bootstrap(context.Background()); err != nil {
		return nil, err
	}

	return adapter, nil
}

// Bootstrap creates the necessary table if it doesn't exist
func (a *LibSQLAdapter) Bootstrap(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS snippets (
		id TEXT PRIMARY KEY,
		title TEXT,
		owner_id TEXT,
		created_at DATETIME
	);`
	_, err := a.db.ExecContext(ctx, query)
	return err
}

// CreateSnippet satisfies the core.IdentityRepository interface
func (a *LibSQLAdapter) CreateSnippet(ctx context.Context, s *core.Snippet) error {
	query := `INSERT INTO snippets (id, title, owner_id, created_at) VALUES (?, ?, ?, ?)`
	_, err := a.db.ExecContext(ctx, query, s.ID, s.Title, s.OwnerID, s.CreatedAt)
	return err
}