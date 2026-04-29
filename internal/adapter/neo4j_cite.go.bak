package adapter

import (
	"context"
	"time"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CiteSnippet implements the University Citation Pattern: Relationship as an Entity
func (a *Neo4jAdapter) CiteSnippet(ctx context.Context, c *core.Citation) error {
	session := a.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (source:Snippet {id: $sourceID})
			MERGE (target:Snippet {id: $targetID})
			CREATE (source)-[r:CITES {
				context: $context,
				created_at: datetime($createdAt)
			}]->(target)
			RETURN r
		`
		params := map[string]any{
			"sourceID":  c.SourceID,
			"targetID":  c.TargetID,
			"context":   c.Context,
			"createdAt": c.Timestamp.Format(time.RFC3339),
		}

		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if !res.Next(ctx) {
			return nil, nil
		}
		return nil, nil
	})

	return err
}