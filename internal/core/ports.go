package core

import "context"

// IdentityRepository defines the 'Shore' behavior (libSQL)
type IdentityRepository interface {
	CreateSnippet(ctx context.Context, s *Snippet) error
	Search(ctx context.Context, query string) ([]*Snippet, error)
	Close() error // Close the repository and release any resources
}

// GraphRepository defines the 'Current' behavior (Neo4j)
type GraphRepository interface {
	//SaveVersion(ctx context.Context, snippetID string, v *Version) error
	SaveVersion(ctx context.Context, snippet *Snippet, v *Version) error
	CiteSnippet(ctx context.Context, c *Citation) error
	Close() error // Close the repository and release any resources
}


