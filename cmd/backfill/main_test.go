package main

import (
	"context"
	"testing"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

// --- 1. Hand-rolled Mocks ---

type mockIdentityRepo struct {
	snippets []*core.Snippet
	err      error // defined as 'err'
}

func (m *mockIdentityRepo) GetAll(ctx context.Context) ([]*core.Snippet, error) {
	return m.snippets, m.err // returning 'err'
}

type mockGraphRepo struct {
	savedSnippets []*core.Snippet
	err           error
}

func (m *mockGraphRepo) SaveVersion(ctx context.Context, snippet *core.Snippet) error {
	m.savedSnippets = append(m.savedSnippets, snippet)
	return m.err
}

// --- 2. The Test ---

func TestBackfillRunner_Success(t *testing.T) {
	// Arrange: Removed the 'Content' field to match your domain struct
	legacySnippets := []*core.Snippet{
		{ID: "snip-1", Title: "Legacy 1"},
		{ID: "snip-2", Title: "Legacy 2"},
	}

	idRepo := &mockIdentityRepo{snippets: legacySnippets}
	graphRepo := &mockGraphRepo{}

	// Instantiate our Runner (this is what will fail now!)
	runner := NewRunner(idRepo, graphRepo)

	// Act: Run the backfill
	err := runner.Run(context.Background())

	// Assert: Verify the results
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(graphRepo.savedSnippets) != 2 {
		t.Fatalf("expected 2 snippets synced to graph, got %d", len(graphRepo.savedSnippets))
	}

	if graphRepo.savedSnippets[0].ID != "snip-1" {
		t.Errorf("expected first synced snippet to be snip-1, got %s", graphRepo.savedSnippets[0].ID)
	}
}