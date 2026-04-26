package http

import (
	"context"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

// SnippetService defines the contract the handler needs (The Port)
type SnippetService interface {
    CreateSnippet(ctx context.Context, title, ownerID, content string) (*core.Snippet, error)
    SearchSnippets(ctx context.Context, query string) ([]*core.Snippet, error)
    CiteSnippet(ctx context.Context, sourceID, targetID, contextStr string) error
}
type Handler struct {
    svc SnippetService // We use 'svc' consistently here
}

func NewHandler(svc SnippetService) *Handler {
    return &Handler{svc: svc}
}