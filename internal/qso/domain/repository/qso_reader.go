package repository

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/qso/domain/model"
)

// QSOReader defines read operations for QSOs
type QSOReader interface {
	// FindByID finds a QSO by ID
	FindByID(ctx context.Context, id string) (*model.QSO, error)

	// FindByUserID finds QSOs for a user with filtering and pagination
	FindByUserID(ctx context.Context, userID string, filter *model.Filter) (*model.QSOLists, error)

	// Search searches QSOs by callsign, notes, or other criteria
	Search(ctx context.Context, userID string, query string, limit int) ([]*model.QSO, error)

	// CountByUserID counts total QSOs for a user
	CountByUserID(ctx context.Context, userID string) (int, error)

	// GetRecentQSOs gets recent QSOs for a user
	GetRecentQSOs(ctx context.Context, userID string, limit int) ([]*model.QSO, error)

	// GetUniqueCallsigns gets unique callsigns for a user
	GetUniqueCallsigns(ctx context.Context, userID string) ([]string, error)
}
