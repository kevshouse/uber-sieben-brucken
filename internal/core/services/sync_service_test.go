package services_test

import (
	"context"
	"testing"

	"github.com/kevshouse/uber-sieben-brucken/internal/core/services"
)

// Mocking the shores for our unit test
type mockIdShore struct{ called bool }
func (m *mockIdShore) Save(ctx context.Context, d any) error { m.called = true; return nil }

type mockGraphShore struct{ called bool }
func (m *mockGraphShore) SyncNode(ctx context.Context, d any) error { m.called = true; return nil }

func TestSyncService_LiveSync(t *testing.T) {
	id := &mockIdShore{}
	graph := &mockGraphShore{}
	srv := services.NewSyncService(id, graph)

	err := srv.LiveSync(context.Background(), "Seven Bridges Data")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !id.called || !graph.called {
		t.Errorf("Sync failed: ID called: %v, Graph called: %v", id.called, graph.called)
	}
}