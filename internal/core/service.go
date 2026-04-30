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

	if err := s.identityRepo.Save(ctx, snippet); err != nil {
		return nil, fmt.Errorf("failed to anchor identity: %w", err)
	}

	v := &Version{
		ID:        uuid.New().String(),
		Content:   content,
		Timestamp: now,
	}

	if err := s.graphRepo.SaveVersion(ctx, snippet, v); err != nil { // Pass 'snippet' here
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

func (s *SnippetService) SearchSnippets(ctx context.Context, query string) ([]*Snippet, error) {
	if query == "" {
		return []*Snippet{}, nil
	}
	results, err := s.identityRepo.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	return results, nil
}
