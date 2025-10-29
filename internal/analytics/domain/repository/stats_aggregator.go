package repository

import (
	"context"
)

// StatsAggregator defines statistical aggregation operations
type StatsAggregator interface {
	// GetQSOsByBand gets QSO count grouped by band
	GetQSOsByBand(ctx context.Context, userID string, dateFrom, dateTo string) (map[string]int, error)

	// GetQSOsByMode gets QSO count grouped by mode
	GetQSOsByMode(ctx context.Context, userID string, dateFrom, dateTo string) (map[string]int, error)

	// GetQSOsByContinent gets QSO count grouped by continent
	GetQSOsByContinent(ctx context.Context, userID string, dateFrom, dateTo string) (map[string]int, error)

	// GetDailyActivity gets daily QSO counts for a period
	GetDailyActivity(ctx context.Context, userID string, days int) ([]int, error)
}
