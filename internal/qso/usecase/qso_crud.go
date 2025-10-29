package usecase

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/qso/domain/model"
)

// QSOCrud defines basic CRUD operations for QSOs
type QSOCrud interface {
	// CreateQSO creates a new QSO record
	CreateQSO(ctx context.Context, req *model.QSO) (*model.QSO, error)

	// GetQSO retrieves a QSO by ID
	GetQSO(ctx context.Context, id, userID string) (*model.QSO, error)

	// UpdateQSO updates an existing QSO
	UpdateQSO(ctx context.Context, qso *model.QSO) (*model.QSO, error)

	// DeleteQSO deletes a QSO by ID
	DeleteQSO(ctx context.Context, id, userID string) error
}
