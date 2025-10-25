package domain

import "time"

// PeriodType представляет тип периода для аналитики
type PeriodType string

const (
	Period1Day   PeriodType = "1d"
	Period7Days  PeriodType = "7d"
	Period30Days PeriodType = "30d"
	Period1Year  PeriodType = "1y"
	PeriodCustom PeriodType = "custom"
)

// TimeRange представляет временной диапазон
type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// ContinentStatistics содержит статистику по континентам
type ContinentStatistics struct {
	Continent  string  `json:"continent"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// GrowthMetrics содержит метрики роста
type GrowthMetrics struct {
	TotalQSOGrowth float64 `json:"total_qso_growth"`
	NewCountries   int     `json:"new_countries"`
	NewGrids       int     `json:"new_grids"`
	ActiveDays     int     `json:"active_days"`
}

// ExportOptions содержит параметры экспорта данных
type ExportOptions struct {
	Format     string    `json:"format"`
	IncludeQSL bool      `json:"include_qsl"`
	DateFrom   time.Time `json:"date_from"`
	DateTo     time.Time `json:"date_to"`
}
