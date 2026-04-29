package adapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type LibSQLAdapter struct {
	db *sql.DB
}

func NewLibSQLAdapter(url string) (*LibSQLAdapter, error) {
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	adapter := &LibSQLAdapter{db: db}
	return adapter, nil
}

// Save now matches the strict core.IdentityShore interface
func (a *LibSQLAdapter) Save(ctx context.Context, s *core.Snippet) error {
    // We no longer need to check if 'entity' is a snippet; 
    // the compiler guarantees it is now.
    
    query := `INSERT INTO snippets (id, title, owner_id, created_at) 
              VALUES (?, ?, ?, ?)
              ON CONFLICT(id) DO UPDATE SET 
              title=excluded.title, owner_id=excluded.owner_id`
    
    _, err := a.db.ExecContext(ctx, query, s.ID, s.Title, s.OwnerID, s.CreatedAt)
    return err
}

// CreateSnippet now matches the strict core.IdentityShore interface
func (a *LibSQLAdapter) CreateSnippet(ctx context.Context, s *core.Snippet) error {
    // No more type assertion needed! 's' is already a *core.Snippet.
    query := `INSERT INTO snippets (id, title, owner_id, created_at) VALUES (?, ?, ?, ?)`
    
    _, err := a.db.ExecContext(ctx, query, s.ID, s.Title, s.OwnerID, s.CreatedAt)
    return err
}

func (a *LibSQLAdapter) Search(ctx context.Context, query string) ([]*core.Snippet, error) {
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
	return snippets, nil
}

func (a *LibSQLAdapter) GetAll(ctx context.Context) ([]*core.Snippet, error) {
    var snippets []*core.Snippet
    query := `SELECT id, title, owner_id, created_at FROM snippets` // Adjust columns to your schema
    
    rows, err := a.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to query snippets: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        s := &core.Snippet{}
        if err := rows.Scan(&s.ID, &s.Title, &s.OwnerID, &s.CreatedAt); err != nil {
            return nil, err
        }
        snippets = append(snippets, s)
    }
    return snippets, nil
}

func (a *LibSQLAdapter) Close() error { return a.db.Close() }