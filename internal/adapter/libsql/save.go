package libsql

import (
	"context"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

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