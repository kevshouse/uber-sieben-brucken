package adapter

//import (
//	"context"
//	"database/sql"
//
//	"github.com/kevshouse/uber-sieben-brucken/internal/core"
//	_ "github.com/ncruces/go-sqlite3/driver" // '_' forces the compiler to keep package for registering 'libsql' at runtime
//)

import (
	"context"
	"database/sql"
	"time"

	// You'll need this for the CreatedAt timestamps
	"github.com/kevshouse/uber-sieben-brucken/internal/core"

	// Switch this back to the libSQL driver we verified yesterday
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)
var _ = time.Now() // This is a "dummy" use to ensure the 'time' package is imported, which is needed for core.Snippet.CreatedAt

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

// The "Looking Glass" (Search & Retrieval)
func (a *LibSQLAdapter) Search(ctx context.Context, query string) ([]*core.Snippet, error) {
	// Using %query% to find the string in the title.
	rows, err := a.db.QueryContext(ctx, `SELECT id, title, owner_id, created_at FROM snippets WHERE title LIKE ?`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*core.Snippet

	for rows.Next() {
		s := &core.Snippet{}
		if err := rows.Scan(&s.ID, &s.Title, &s.OwnerID, &s.CreatedAt); err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if snippets == nil {
		return []*core.Snippet{}, nil // Return an empty slice instead of nil
	}
	return snippets, nil
}