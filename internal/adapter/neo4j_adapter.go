package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jAdapter struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jAdapter(uri, username, password string) (*Neo4jAdapter, error) {
    // The driver is the engine under the hood
    driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
    if err != nil {
        return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
    }

    // Return the pointer to our adapter struct
    return &Neo4jAdapter{
        driver: driver,
    }, nil
}

func (a *Neo4jAdapter) SyncNode(ctx context.Context, snippet *core.Snippet) error {
	return nil // Placeholder for logic
}


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

// CiteSnippet now matches the strict core.GraphShore interface
func (a *Neo4jAdapter) CiteSnippet(ctx context.Context, c *core.Citation) error {
    // Note: We removed 'entity any' and the 'if !ok' block entirely!
    session := a.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer session.Close(ctx)

    _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
        query := `MATCH (s:Snippet {id: $sourceID}) 
                  MERGE (t:Snippet {id: $targetID}) 
                  CREATE (s)-[r:CITES {context: $ctx}]->(t) 
                  RETURN r`
        return tx.Run(ctx, query, map[string]any{
            "sourceID": c.SourceID,
            "targetID": c.TargetID,
            "ctx":      c.Context,
        })
    })
    return err
}

func (a *Neo4jAdapter) Close() error { return a.driver.Close(context.Background()) }