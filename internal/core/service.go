package core

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SnippetService struct {
	identityRepo IdentityRepository
	graphRepo    GraphRepository
}

func NewSnippetService(idRepo IdentityRepository, gRepo GraphRepository) *SnippetService {
	return &SnippetService{
		identityRepo: idRepo,
		graphRepo:    gRepo,
	}
}

func (s *SnippetService) CreateSnippet(ctx context.Context, title, ownerID, content string) (*Snippet, error) {
	id := uuid.New().String()
	now := time.Now()

	snippet := &Snippet{
		ID:        id,
		Title:     title,
		OwnerID:   ownerID,
		CreatedAt: now,
	}

	if err := s.identityRepo.CreateSnippet(ctx, snippet); err != nil {
		return nil, fmt.Errorf("failed to anchor identity: %w", err)
	}

	v := &Version{
		ID:        uuid.New().String(),
		Content:   content,
		Timestamp: now,
	}

	if err := s.graphRepo.SaveVersion(ctx, id, v); err != nil {
		return nil, fmt.Errorf("failed to anchor graph lineage: %w", err)
	}

	return snippet, nil
}

func (s *SnippetService) CiteSnippet(ctx context.Context, sourceID, targetID, contextStr string) error {
	citation := &Citation{
		ID:        uuid.New().String(),
		SourceID:  sourceID,
		TargetID:  targetID,
		Context:   contextStr,
		Timestamp: time.Now(),
	}

	return s.graphRepo.CiteSnippet(ctx, citation)
}