package repository

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/qso/domain/model"
)

// QSOWriter defines write operations for QSOs
type QSOWriter interface {
	// Create creates a new QSO
	Create(ctx context.Context, qso *model.QSO) error

	// Update updates an existing QSO
	Update(ctx context.Context, qso *model.QSO) error

	// Delete deletes a QSO by ID
	Delete(ctx context.Context, id string) error
}
