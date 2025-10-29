package usecase

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/auth/domain/model"
)

// ProfileManager defines user profile operations
type ProfileManager interface {
	// UpdateProfile updates user profile information
	UpdateProfile(ctx context.Context, userID string, req model.UpdateProfileRequest) (*model.User, error)

	// GetProfile retrieves user profile information
	GetProfile(ctx context.Context, userID string) (*model.User, error)

	// ChangePassword changes user password
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
}
