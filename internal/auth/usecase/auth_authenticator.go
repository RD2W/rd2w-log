package usecase

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/auth/domain/model"
)

// AuthAuthenticator defines authentication operations
type AuthAuthenticator interface {
	// Register registers a new user
	Register(ctx context.Context, req model.RegisterRequest) (*model.User, *model.TokenPair, error)

	// Login authenticates a user and returns tokens
	Login(ctx context.Context, creds model.LoginCredentials) (*model.User, *model.TokenPair, error)

	// ValidateToken validates an access token and returns user information
	ValidateToken(ctx context.Context, token string) (*model.User, error)

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error)
}
