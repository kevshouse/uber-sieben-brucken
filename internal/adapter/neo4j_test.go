package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j" // Required for the helper
	"github.com/stretchr/testify/assert"        // This makes checks much cleaner
)

func TestNeo4jAdapter(t *testing.T) {
	ctx := context.Background()

	// 1. Initialize the Adapter
	repo, err := NewNeo4jAdapter("bolt://localhost:7687", "neo4j", "password123")
	assert.NoError(t, err)

	// 2. Database Cleanup (Your "Clean Shore" Policy)
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, "MATCH (n) DETACH DELETE n", nil)
	})
	session.Close(ctx)
	assert.NoError(t, err)

	// --- SUB-TEST 1: Versions ---
	t.Run("Version Evolution", func(t *testing.T) {
		snippetID := "test-snippet-001"
		
		// Crucial: We must create the Snippet anchor before adding versions
		err = setupSnippetAnchor(ctx, repo, snippetID)
		assert.NoError(t, err)

		v1 := &core.Version{ID: "v1", Content: "fmt.Println('Genesis')", Timestamp: time.Now()}
		err = repo.SaveVersion(ctx, snippetID, v1)
		assert.NoError(t, err)

		v2 := &core.Version{ID: "v2", Content: "fmt.Println('Evolution')", Timestamp: time.Now().Add(time.Minute)}
		err = repo.SaveVersion(ctx, snippetID, v2)
		assert.NoError(t, err)

		history, err := repo.GetHistory(ctx, snippetID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(history))
	})

	// --- SUB-TEST 2: University Citation Pattern ---
	t.Run("University Citation Pattern", func(t *testing.T) {
		sourceID := "snippet-a"
		targetID := "snippet-b"

		// Create the two shores
		assert.NoError(t, setupSnippetAnchor(ctx, repo, sourceID))
		assert.NoError(t, setupSnippetAnchor(ctx, repo, targetID))

		// Create first citation
		cit1 := &core.Citation{
			ID:        "cit-1",
			SourceID:  sourceID,
			TargetID:  targetID,
			Context:   "Initial reference",
			Timestamp: time.Now(),
		}
		err = repo.CiteSnippet(ctx, cit1)
		assert.NoError(t, err)

		// Evolve the citation (The "Relationship as Entity" check)
		cit2 := &core.Citation{
			ID:        "cit-2",
			SourceID:  sourceID,
			TargetID:  targetID,
			Context:   "Updated context",
			Timestamp: time.Now().Add(time.Minute),
		}
		err = repo.CiteSnippet(ctx, cit2)
		assert.NoError(t, err)
	})
}

// Helper: This bridges the gap between libSQL and Neo4j
func setupSnippetAnchor(ctx context.Context, a *Neo4jAdapter, id string) error {
	session := a.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, "MERGE (s:Snippet {id: $id})", map[string]interface{}{"id": id})
		return nil, err
	})
	return err
}