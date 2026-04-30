package libsql

import (
	"context"
	"fmt"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)


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
