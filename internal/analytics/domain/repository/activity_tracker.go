package repository

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/analytics/domain/model"
)

// ActivityTracker defines activity tracking operations
type ActivityTracker interface {
	// GetActivityHeatmap gets activity data for heatmap visualization
	GetActivityHeatmap(ctx context.Context, userID, period string) ([]*model.ActivityPoint, error)
}
