package usecase

import "context"

// SessionCleaner defines session cleanup operations
type SessionCleaner interface {
	// Logout invalidates user sessions
	Logout(ctx context.Context, userID string) error
}
