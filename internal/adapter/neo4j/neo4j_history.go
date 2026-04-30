package neo4j

import (
	"context"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (a *Neo4jAdapter) GetHistory(ctx context.Context, snippetID string) ([]*core.Version, error) {
	session := a.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (n:Snippet {id: $id})-[:HAS_VERSION]->(v:Version)
			RETURN v.id as id, v.content as content
			ORDER BY v.created_at DESC
		`
		res, err := tx.Run(ctx, query, map[string]any{"id": snippetID})
		if err != nil {
			return nil, err
		}

		var versions []*core.Version
		for res.Next(ctx) {
			record := res.Record()
			id, _ := record.Get("id")
			content, _ := record.Get("content")

			versions = append(versions, &core.Version{
				ID:      id.(string),
				Content: content.(string),
			})
		}
		return versions, nil
	})

	if err != nil {
		return nil, err
	}
	return result.([]*core.Version), nil
}