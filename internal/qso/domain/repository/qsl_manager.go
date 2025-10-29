package repository

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/qso/domain/model"
)

// QSLManager defines QSL management operations
type QSLManager interface {
	// Create creates QSL information for a QSO
	Create(ctx context.Context, qsl *model.QSLInfo) error

	// FindByQSOID finds QSL information by QSO ID
	FindByQSOID(ctx context.Context, qsoID string) (*model.QSLInfo, error)

	// Update updates QSL information
	Update(ctx context.Context, qsl *model.QSLInfo) error
}
