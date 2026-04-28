package ports

import "context"

// Identity ddefines how we persist the source of truth
type IdentityShore interface {
	Save(ctx context.Context, entity any) error
}

// GraphShore defines how we persist relational connections.
type GraphShore interface {
	SyncNode(ctx context.Context, entity any) error
}
