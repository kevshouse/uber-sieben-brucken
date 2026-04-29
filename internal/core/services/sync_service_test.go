package services_test

import (
	"context"
	"testing"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
	"github.com/kevshouse/uber-sieben-brucken/internal/core/services"
)

// Mocking the shores for our unit test.
type mockIdShore struct {
	called bool
}

// Updated to match core.IdentityShore signatures exactly
func (m *mockIdShore) Save(ctx context.Context, s *core.Snippet) error {
	m.called = true
	return nil
}
func (m *mockIdShore) CreateSnippet(ctx context.Context, s *core.Snippet) error { return nil }
func (m *mockIdShore) Search(ctx context.Context, q string) ([]*core.Snippet, error) { return nil, nil }
func (m *mockIdShore) GetAll(ctx context.Context) ([]*core.Snippet, error)             { return nil, nil }
func (m *mockIdShore) Close() error                                                  { return nil }

type mockGraphShore struct {
	called bool
}

// Updated to match core.GraphShore signatures exactly
func (m *mockGraphShore) SyncNode(ctx context.Context, s *core.Snippet) error {
	m.called = true
	return nil
}
func (m *mockGraphShore) SaveVersion(ctx context.Context, s *core.Snippet, v *core.Version) error {
	return nil
}
func (m *mockGraphShore) CiteSnippet(ctx context.Context, c *core.Citation) error { return nil }
func (m *mockGraphShore) Close() error                                          { return nil }

func TestSyncService_LiveSync(t *testing.T) {
	t.Run("it should successfully persist to both shores", func(t *testing.T) {
		id := &mockIdShore{}
		graph := &mockGraphShore{}
		srv := services.NewSyncService(id, graph)

		// Create a concrete *core.Snippet instead of passing a string
		testSnippet := &core.Snippet{
			ID:    "test-id",
			Title: "Seven Bridges Data",
		}

		err := srv.LiveSync(context.Background(), testSnippet)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !id.called {
			t.Error("Identity Shore (libSQL) was not reached")
		}

		if !graph.called {
			t.Error("Graph Shore (Neo4j) was not reached")
		}
	})
}