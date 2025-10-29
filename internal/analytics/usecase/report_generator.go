package usecase

import (
	"context"

	"github.com/rd2w/rd2w-log/internal/analytics/domain/model"
)

// ReportGenerator defines report generation operations
type ReportGenerator interface {
	// GenerateReport generates a report in the specified format
	GenerateReport(ctx context.Context, req *model.ReportRequest) (*model.ReportResponse, error)
}
