package usecase

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/qso/domain/model"
)

// QSLHandler defines QSL card management operations
type QSLHandler interface {
	// UpdateQSLInfo updates QSL information for a QSO
	UpdateQSLInfo(ctx context.Context, qsoID string, qsl *model.QSLInfo) error

	// GetQSLInfo retrieves QSL information for a QSO
	GetQSLInfo(ctx context.Context, qsoID string) (*model.QSLInfo, error)
}
