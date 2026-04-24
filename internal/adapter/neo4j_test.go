package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/assert"
)

func TestNeo4jAdapter(t *testing.T) {
	ctx := context.Background()

	// 1. Initialize the Adapter
	repo, err := NewNeo4jAdapter("bolt://localhost:7687", "neo4j", "password123")
	assert.NoError(t, err)

	// 2. Database Cleanup
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, "MATCH (n) DETACH DELETE n", nil)
	})
	session.Close(ctx)
	assert.NoError(t, err)

	// --- SUB-TEST 1: Versions ---
	t.Run("Version Evolution", func(t *testing.T) {
		snippetID := "test-snippet-001"
		
		// Create the Snippet struct to satisfy the new Port signature
		s := &core.Snippet{
			ID:    snippetID,
			Title: "Test Snippet",
		}

		v1 := &core.Version{
			ID:      "v1",
			Content: "First version",
		}
		v2 := &core.Version{
			ID:      "v2",
			Content: "Second version",
		}

		// Initial creation
		err = repo.SaveVersion(ctx, s, v1)
		assert.NoError(t, err)

		// Evolution
		err = repo.SaveVersion(ctx, s, v2)
		assert.NoError(t, err)

		history, err := repo.GetHistory(ctx, snippetID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(history))
	})

	// --- SUB-TEST 2: University Citation Pattern ---
	t.Run("University Citation Pattern", func(t *testing.T) {
		sourceID := "snippet-a"
		targetID := "snippet-b"

		// Create anchors using the updated helper
		assert.NoError(t, setupSnippetAnchor(ctx, repo, sourceID))
		assert.NoError(t, setupSnippetAnchor(ctx, repo, targetID))

		cit1 := &core.Citation{
			ID:        "cit-1",
			SourceID:  sourceID,
			TargetID:  targetID,
			Context:   "Initial reference",
			Timestamp: time.Now(),
		}
		err = repo.CiteSnippet(ctx, cit1)
		assert.NoError(t, err)
	})
}

// Helper updated to match the new Port signature
func setupSnippetAnchor(ctx context.Context, repo *Neo4jAdapter, snippetID string) error {
	s := &core.Snippet{
		ID:    snippetID,
		Title: "Anchor Snippet",
	}
	v := &core.Version{
		ID:      "anchor-v-" + snippetID,
		Content: "Anchor content",
	}
	return repo.SaveVersion(ctx, s, v)
}