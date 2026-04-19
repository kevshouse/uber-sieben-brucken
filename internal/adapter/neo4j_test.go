package adapter

import (
	"context"
	"testing"
	"time"
	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

func TestNeo4jAdapter_SaveVersion(t *testing.T) {
	ctx := context.Background()
	
	// 1. Setup our Adapter (this will fail to compile initially)
	// We'll assume a constructor NewNeo4jAdapter exists later
	repo := NewNeo4jAdapter("bolt://localhost:7687", "neo4j", "password123")
	
	snippetID := "test-snippet-001"
	v1 := &core.Version{
		ID:        "v1",
		Content:   "fmt.Println('Hello, Bridge 1')",
		Timestamp: time.Now(),
	}

	// 2. The First Act: Save Genesis Version
	err := repo.SaveVersion(ctx, snippetID, v1)
	if err != nil {
		t.Fatalf("Failed to save genesis version: %v", err)
	}

	// 3. The Second Act: Save Version 2 (The Handshake)
	v2 := &core.Version{
		ID:        "v2",
		Content:   "fmt.Println('Hello, Bridge 2')",
		Timestamp: time.Now(),
	}
	
	err = repo.SaveVersion(ctx, snippetID, v2)
	if err != nil {
		t.Fatalf("Failed to save version 2: %v", err)
	}

	// 4. The Verification: This is where we test the Spotnet Spec
	// We will write a test-only query to verify the bidirectional edges
	history, err := repo.GetHistory(ctx, snippetID)
	if err != nil {
		t.Fatalf("Failed to retrieve history: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("Expected 2 versions in history, got %d", len(history))
	}
    
    // Check for bidirectional link between V2 and V1
    // Logic: history[0] should be V2, history[1] should be V1
}