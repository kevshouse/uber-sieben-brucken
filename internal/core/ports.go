package core

import "context"

// IdentityRepository defines the 'Shore' behavior (libSQL)
type IdentityRepository interface {
	CreateSnippet(ctx context.Context, s *Snippet) error
}

// GraphRepository defines the 'Current' behavior (Neo4j)
type GraphRepository interface {
	SaveVersion(ctx context.Context, snippetID string, v *Version) error
	CiteSnippet(ctx context.Context, c *Citation) error
}