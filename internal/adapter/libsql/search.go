package libsql

import (
	"context"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

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