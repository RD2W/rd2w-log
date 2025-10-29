package usecase

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/analytics/domain/model"
)

// ActivityVisualizer defines activity visualization operations
type ActivityVisualizer interface {
	// GetActivityHeatmap gets activity heatmap data
	GetActivityHeatmap(ctx context.Context, userID, period string) (*model.ActivityHeatmap, error)
}
