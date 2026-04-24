package adapter

import (
	"context"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jAdapter struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jAdapter(uri, username, password string) (*Neo4jAdapter, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	return &Neo4jAdapter{driver: driver}, nil
}

func (a *Neo4jAdapter) SaveVersion(ctx context.Context, s *core.Snippet, v *core.Version) error {
	session := a.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
		MERGE (n:Snippet {id: $snippetID})
		SET n.title = $title, n.owner_id = $ownerID
		WITH n  // <--- This is the essential bridge!
		OPTIONAL MATCH (n)-[oldRel:LATEST_VERSION]->(oldV:Version)
		CREATE (newV:Version {id: $id, content: $content, ts: $ts})
		CREATE (n)-[:LATEST_VERSION]->(newV)
		DELETE oldRel
		WITH oldV, newV
		WHERE oldV IS NOT NULL
		CREATE (newV)-[:PREVIOUS]->(oldV)
		RETURN newV.id
	`
		params := map[string]interface{}{
			"snippetID": s.ID,
			"title"	:     s.Title,
			"ownerID":   s.OwnerID,
			"id":        v.ID,
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
		query := `
			MATCH (s:Snippet {id: $snippetID})-[:LATEST_VERSION]->(v:Version)
			MATCH (v)-[:PREVIOUS*0..]->(history)
			RETURN history.id, history.content, history.ts ORDER BY history.ts DESC
		`
		res, err := tx.Run(ctx, query, map[string]interface{}{"snippetID": snippetID})
		if err != nil {
			return nil, err
		}

		var versions []*core.Version
		for res.Next(ctx) {
			record := res.Record()
			id, _ := record.Get("history.id")
			content, _ := record.Get("history.content")
			_, _ = record.Get("history.ts")

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

// CiteSnippet implements the University Citation Pattern: Relationship as an Entity
func (a *Neo4jAdapter) CiteSnippet(ctx context.Context, c *core.Citation) error {
	session := a.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (source:Snippet {id: $sourceID})
			MATCH (target:Snippet {id: $targetID})
			
			OPTIONAL MATCH (source)-[oldRel:CITES_LATEST]->(prevCit:Citation)-[:CITES_TARGET]->(target)
			
			CREATE (newCit:Citation {
				id: $citID, 
				context: $context, 
				ts: $ts
			})
			
			CREATE (source)-[:CITES_LATEST]->(newCit)
			CREATE (newCit)-[:CITES_TARGET]->(target)
			
			DELETE oldRel
			WITH prevCit, newCit
			WHERE prevCit IS NOT NULL
			CREATE (newCit)-[:PREVIOUS]->(prevCit)
			CREATE (prevCit)-[:NEXT]->(newCit)
			
			RETURN newCit.id
		`
		params := map[string]interface{}{
			"sourceID": c.SourceID,
			"targetID": c.TargetID,
			"citID":    c.ID,
			"context":  c.Context,
			"ts":       c.Timestamp.Unix(),
		}
		return tx.Run(ctx, query, params)
	})
	return err
}