package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

type mockIdShore struct{
	called bool
}

func (m *mockIdShore) Save(ctx context.Context, s *core.Snippet) error {
	m.called = true
	return nil
}

type SyncService struct {
	identity core.IdentityShore
	graph    core.GraphShore
}

func NewSyncService(id core.IdentityShore, g core.GraphShore) *SyncService {
	return &SyncService{
		identity: id,
		graph:    g,
	}
}

func (s *SyncService) Sync(ctx context.Context, snippet *core.Snippet) error {
	return s.graph.SyncNode(ctx, snippet)
}

// LiveSync coordinates the synchronous transfer of data to both shores, ensuring consistency and reliability.
func (s *SyncService) LiveSync(ctx context.Context, snippet *core.Snippet) error {
	// The Core Domain is the master of time.
	if snippet.CreatedAt.IsZero() {
		snippet.CreatedAt = time.Now()
	}

	// 1. Persist to Identity Shore
	if err := s.identity.Save(ctx, snippet); err != nil {
		return fmt.Errorf("failed to save to identity shore: %w", err)
	}

	// 2. Persist to Graph Shore
	if err := s.graph.SyncNode(ctx, snippet); err != nil {
		// Since we are synchrounous, we reporrt this immeadiately.
		// In a later refactor, we might want to implement a compensation
		// mechanism here to rollback or retry the identity shore if the graph shore fails.
		return fmt.Errorf("failed to sync to graph shore: %w", err)
	}

	return nil
}

// CreateGenesis acts as the secure orchestrator, buiding both the static identity
// and the first immutable version of the snippet's history
func (s *SyncService) CreateGenesis(ctx context.Context, title, ownerID string, content string) (*core.Snippet, error) {
	// 1. The Core Domain establishes absolute truth
	now := time.Now()
	snippetID := uuid.New().String()
	versionID := uuid.New().String()

	// 2. Build the Identity (The Folder)
	snippet := &core.Snippet{
		ID:        snippetID,
		Title:     title,
		OwnerID:   ownerID,
		CreatedAt: now,
	}

	// 3. Build the Path (The First Page)
	version := &core.Version{
		ID:        versionID,
		Content:   content,
		Timestamp: now,
	}

	// 4. Persist to Identity Shore
	if err := s.identity.Save(ctx, snippet); err != nil {
		return nil, fmt.Errorf("failed to save genesis to libSQL: %w", err)
	}

	// 5. Persist to Graph Shore
	if err := s.graph.SaveVersion(ctx, snippet, version); err != nil {
		// Note for our technical debt log: If this fails, we have a "Ghost Shore"
		// identity without history. We will add a compensation/rollback mechanism later!
		return nil, fmt.Errorf("failed to save genesis version to Neo4j: %w", err)
	}

	return snippet, nil

}
