package http

import (
	"context"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

// SnippetService defines the contract the handler needs (The Port).
// The Driving Adapter must conform to this interface.
type SnippetService interface {
	// CreateGenesis acts as the secure orchestrator, taking the raw intent
	// and building both the Snippet (Identity) and the first Version (Path).
	CreateGenesis(ctx context.Context, title, ownerID, content string) (*core.Snippet, error)
	
	SearchSnippets(ctx context.Context, query string) ([]*core.Snippet, error)
	CiteSnippet(ctx context.Context, sourceID, targetID, contextStr string) error
}

type Handler struct {
	svc SnippetService // We use 'svc' consistently here to point to our Brain
}

func NewHandler(svc SnippetService) *Handler {
	return &Handler{svc: svc}
}