package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestNeo4jAdapter_SaveVersion(t *testing.T) {
	ctx := context.Background()

	// 1. Initialize the Adapter
	// Fixed: NewNeo4jAdapter returns two values (*Adapter, error)
	repo, err := NewNeo4jAdapter("bolt://localhost:7687", "neo4j", "password123")
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	// 2. Database Cleanup (The "Clean Shore" Policy)
	// We wipe the DB before starting to avoid the "Expected 2, Got 4" pothole.
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, "MATCH (n) DETACH DELETE n", nil)
	})
	session.Close(ctx)
	if err != nil {
		t.Fatalf("Failed to clean database: %v", err)
	}

	snippetID := "test-snippet-001"

	// 3. Save Genesis Version (V1)
	v1 := &core.Version{
		ID:        "v1",
		Content:   "fmt.Println('Genesis Bridge')",
		Timestamp: time.Now(),
	}

	// Fixed: Use '=' instead of ':=' because 'err' was already declared at line 17
	err = repo.SaveVersion(ctx, snippetID, v1)
	if err != nil {
		t.Fatalf("Failed to save genesis version: %v", err)
	}

	// 4. Save Second Version (V2 - The Handshake)
	v2 := &core.Version{
		ID:        "v2",
		Content:   "fmt.Println('Second Bridge')",
		Timestamp: time.Now().Add(time.Minute), // Ensure chronological order
	}

	err = repo.SaveVersion(ctx, snippetID, v2)
	if err != nil {
		t.Fatalf("Failed to save version 2: %v", err)
	}

	// 5. Verification
	history, err := repo.GetHistory(ctx, snippetID)
	if err != nil {
		t.Fatalf("Failed to retrieve history: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("Expected 2 versions in history, got %d", len(history))
	}

	// Verify order (The Spotnet Bidirectional Check)
	if history[0].ID != "v1" || history[1].ID != "v2" {
		t.Errorf("History order is incorrect. Got %s -> %s", history[0].ID, history[1].ID)
	}
}