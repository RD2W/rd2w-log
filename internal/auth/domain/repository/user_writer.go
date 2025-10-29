package repository

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/auth/domain/model"
)

// UserWriter defines write operations for users
type UserWriter interface {
	// Create creates a new user
	Create(ctx context.Context, user *model.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *model.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error
}
