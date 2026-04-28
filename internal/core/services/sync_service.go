package services

import (
	"context"
	"fmt"

	"github.com/kevshouse/uber-sieben-brucken/internal/core/ports"
)
	

type SyncService struct {
	identity 	ports.IdentityShore
	graph 		ports.GraphShore
}

func NewSyncService(id ports.IdentityShore, g ports.GraphShore) *SyncService {
	return &SyncService{
		identity: id,
		graph: g,
	}
}

// LiveSync coordinates the synchronous transfer of data to both shores, ensuring consistency and reliability.
func (s *SyncService) LiveSync(ctx context.Context, data any) error {
	// 1. Persist to Identity Shore
	if err := s.identity.Save(ctx, data); err != nil {
		return fmt.Errorf("failed to save to identity shore: %w", err)
	}

	// 2. Persist to Graph Shore
	if err := s.graph.SyncNode(ctx, data); err != nil {
		// Since we are synchrounous, we reporrt this immeadiately.
		// In a later refactor, we might want to implement a compensation 
		// mechanism here to rollback or retry the identity shore if the graph shore fails.
		return fmt.Errorf("failed to sync to graph shore: %w", err)
	}

	return nil
}
