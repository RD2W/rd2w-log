package auth

import "context"

// TokenService defines the interface for token operations
type TokenService interface {
	// GenerateAccessToken generates a new access token
	GenerateAccessToken(userID string) (string, error)

	// GenerateRefreshToken generates a new refresh token
	GenerateRefreshToken(userID string) (string, error)

	// ValidateToken validates a token and returns user ID
	ValidateToken(token string) (string, error)

	// ParseToken parses token claims without validation
	ParseToken(token string) (map[string]interface{}, error)

	// RevokeToken revokes a specific token (для refresh токенов)
	RevokeToken(ctx context.Context, token string) error

	// IsTokenRevoked checks if a token is revoked
	IsTokenRevoked(ctx context.Context, token string) (bool, error)
}
