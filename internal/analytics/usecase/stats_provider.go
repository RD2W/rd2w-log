package usecase

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/analytics/domain/model"
)

// StatsProvider defines statistics retrieval operations
type StatsProvider interface {
	// GetStatistics gets comprehensive statistics for a user
	GetStatistics(ctx context.Context, userID string, dateFrom, dateTo string) (*model.Statistics, error)
}
