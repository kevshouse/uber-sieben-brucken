package core

import "context"

type VersionService interface {
	SaveVersion(ctx context.Context, snippetID string, v *Version) error
	GetHistory(ctx context.Context, snippetID string) ([]*Version, error)

	// THE CITATION BRIDGE: The Citation Handshake
	CiteSnippet(ctx context.Context, citatio *Citation) error
}