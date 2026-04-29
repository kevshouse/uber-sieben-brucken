package adapter

import (
	"context"
	"time"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (a *Neo4jAdapter) SaveVersion(ctx context.Context, s *core.Snippet, v *core.Version) error {
	session := a.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
		MERGE (n:Snippet {id: $snippetID})
		SET n.title = $title, n.owner_id = $ownerID
		WITH n`

		params := map[string]any{
			"snippetID": s.ID,
			"title":     s.Title,
			"ownerID":   s.OwnerID,
		}

		// If a version is provided (standard flow), create it.
		// If v is nil (backfill flow), we just ensure the Snippet anchor exists.
		if v != nil {
			query += `
			OPTIONAL MATCH (n)-[oldRel:LATEST_VERSION]->(oldV:Version)
			CREATE (newV:Version {id: $versionID, content: $content, created_at: datetime($ts)})
			CREATE (n)-[:HAS_VERSION]->(newV)
			DELETE oldRel
			CREATE (n)-[:LATEST_VERSION]->(newV)`
			params["versionID"] = v.ID
			params["content"] = v.Content
			params["ts"] = time.Now().Format(time.RFC3339)
		}

		query += ` RETURN n`
		return tx.Run(ctx, query, params)
	})

	return err
}