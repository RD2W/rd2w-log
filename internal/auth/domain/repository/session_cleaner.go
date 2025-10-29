package repository

import "context"

// SessionCleaner defines session cleanup operations
type SessionCleaner interface {
	// DeleteByUserID deletes all sessions for a user
	DeleteByUserID(ctx context.Context, userID string) error
}
