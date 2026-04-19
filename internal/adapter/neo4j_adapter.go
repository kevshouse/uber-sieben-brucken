package adapter

import (
	"context"
	"time"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jAdapter struct {
	driver neo4j.DriverWithContext
}

// Fixed: Returns (*Neo4jAdapter, error) to match the test's second error
func NewNeo4jAdapter(uri, user, pass string) (*Neo4jAdapter, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, pass, ""))
	if err != nil {
		return nil, err
	}
	return &Neo4jAdapter{driver: driver}, nil
}

func (a *Neo4jAdapter) SaveVersion(ctx context.Context, snippetID string, v *core.Version) error {
	session := a.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// Fixed: Changed tx type to neo4j.ManagedTransaction
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MERGE (anchor:Snippet {id: $snippetID})
			WITH anchor
			OPTIONAL MATCH (anchor)-[oldRel:HAS_LATEST]->(prev:Version)
			CREATE (new:Version {id: $vID, content: $content, ts: $ts})
			CREATE (anchor)-[:HAS_LATEST]->(new)
			DELETE oldRel
			WITH prev, new
			WHERE prev IS NOT NULL
			CREATE (new)-[:PREVIOUS]->(prev)
			CREATE (prev)-[:NEXT]->(new)
			RETURN new.id
		`
		params := map[string]interface{}{
			"snippetID": snippetID,
			"vID":       v.ID,
			"content":   v.Content,
			"ts":        v.Timestamp.Unix(),
		}
		return tx.Run(ctx, query, params)
	})
	return err
}

func (a *Neo4jAdapter) GetHistory(ctx context.Context, snippetID string) ([]*core.Version, error) {
	session := a.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		// This query finds the snippet and follows the chain of versions
		query := `
			MATCH (s:Snippet {id: $snippetID})-[:HAS_LATEST]->(latest:Version)
			MATCH (latest)-[:PREVIOUS*0..]->(v:Version)
			RETURN v.id, v.content, v.ts
			ORDER BY v.ts ASC
		`
		res, err := tx.Run(ctx, query, map[string]interface{}{"snippetID": snippetID})
		if err != nil {
			return nil, err
		}

		var history []*core.Version
		for res.Next(ctx) {
			record := res.Record()
			vID, _ := record.Get("v.id")
			content, _ := record.Get("v.content")
			ts, _ := record.Get("v.ts")

			history = append(history, &core.Version{
				ID:        vID.(string),
				Content:   content.(string),
				Timestamp: time.Unix(ts.(int64), 0),
			})
		}
		return history, nil
	})

	if err != nil {
		return nil, err
	}
	return result.([]*core.Version), nil
}