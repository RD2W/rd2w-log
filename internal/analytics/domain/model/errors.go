package model

import "errors"

var (
	ErrInvalidReportType      = errors.New("invalid report type")
	ErrInvalidReportFormat    = errors.New("invalid report format")
	ErrNoDataForPeriod        = errors.New("no data for specified period")
	ErrReportGenerationFailed = errors.New("report generation failed")
	ErrStatisticsNotFound     = errors.New("statistics not found")
	ErrInvalidPeriod          = errors.New("invalid period")
	ErrInvalidDateRange       = errors.New("invalid date range")
)
