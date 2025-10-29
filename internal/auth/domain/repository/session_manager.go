package repository

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/auth/domain/model"
)

// SessionManager defines session management operations
type SessionManager interface {
	// Create creates a new session
	Create(ctx context.Context, session *model.Session) error

	// FindByAccessToken finds a session by access token
	FindByAccessToken(ctx context.Context, token string) (*model.Session, error)

	// FindByRefreshToken finds a session by refresh token
	FindByRefreshToken(ctx context.Context, token string) (*model.Session, error)

	// Delete deletes a session by ID
	Delete(ctx context.Context, id string) error

	// DeleteExpired deletes all expired sessions
	DeleteExpired(ctx context.Context) error
}
