package migration

import "context"

// SessionAffinity stores the upstream owner of a registration session.
type SessionAffinity interface {
	SetLegacy(ctx context.Context, sessionID string) error
	IsLegacy(ctx context.Context, sessionID string) (bool, error)
}
