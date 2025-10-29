package repository

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/auth/domain/model"
)

// UserReader defines read operations for users
type UserReader interface {
	// FindByID finds a user by ID
	FindByID(ctx context.Context, id string) (*model.User, error)

	// FindByCallsign finds a user by callsign
	FindByCallsign(ctx context.Context, callsign string) (*model.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// ExistsByCallsign checks if a user exists with the given callsign
	ExistsByCallsign(ctx context.Context, callsign string) (bool, error)

	// ExistsByEmail checks if a user exists with the given email
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
