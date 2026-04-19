package core

import "context"

type VersionService interface {
	SaveVersion(ctx context.Context, snippetID string, v *Version) error
	GetHistory(ctx context.Context, snippetID string) ([]*Version, error)
}