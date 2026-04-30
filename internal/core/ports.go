package core

import (
	"context"
)

// IdentityShore defines how we persist the source of truth (libSQL)
type IdentityShore interface {
	Save(ctx context.Context, snippet *Snippet) error
	Search(ctx context.Context, query string) ([]*Snippet, error)
	GetAll(ctx context.Context) ([]*Snippet, error)
	Close() error
}

// GraphShore defines how we persist relational connections (Neo4j)
type GraphShore interface {
	SyncNode(ctx context.Context, snippet *Snippet) error
	SaveVersion(ctx context.Context, snippet *Snippet, v *Version) error
	CiteSnippet(ctx context.Context, c *Citation) error
	Close() error
}

// Aliases to keep the Service and Adapters happy while we finish renaming
type IdentityRepository = IdentityShore
type GraphRepository = GraphShore