package usecase

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/qso/domain/model"
)

// QSOSearcher defines QSO search operations
type QSOSearcher interface {
	// ListQSOs lists QSOs with filtering and pagination
	ListQSOs(ctx context.Context, userID string, filter *model.Filter) (*model.QSOLists, error)

	// SearchQSOs searches QSOs by various criteria
	SearchQSOs(ctx context.Context, userID, query string, limit int) ([]*model.QSO, error)
}
