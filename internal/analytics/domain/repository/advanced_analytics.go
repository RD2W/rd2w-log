package repository

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/analytics/domain/model"
)

// AdvancedAnalytics defines extended analytics operations
type AdvancedAnalytics interface {
	// GetGrowthMetrics gets growth comparison metrics
	GetGrowthMetrics(ctx context.Context, userID string, currentRange, previousRange model.TimeRange) (*model.GrowthMetrics, error)

	// GetContinentBreakdown gets detailed continent statistics
	GetContinentBreakdown(ctx context.Context, userID string, dateFrom, dateTo string) ([]*model.ContinentStatistics, error)

	// GetActivityByPeriod gets activity data for specific periods
	GetActivityByPeriod(ctx context.Context, userID string, period model.PeriodType) (*model.ActivityHeatmap, error)

	// GenerateExport prepares data for export
	GenerateExport(ctx context.Context, userID string, opts model.ExportOptions) ([]byte, error)
}
