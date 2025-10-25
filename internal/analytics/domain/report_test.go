package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rd2w/rd2w-log/internal/analytics/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReportRequest(t *testing.T) {
	userID := uuid.New()
	reportType := domain.ReportTypeDXCC
	format := domain.ReportFormatPDF

	req := domain.NewReportRequest(userID, reportType, format)

	assert.Equal(t, userID, req.UserID)
	assert.Equal(t, reportType, req.Type)
	assert.Equal(t, format, req.Format)
	assert.NotNil(t, req.Options)
	assert.Empty(t, req.Options)
	assert.True(t, req.DateFrom.IsZero())
	assert.True(t, req.DateTo.IsZero())
}

func TestReportRequest_Validate(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		setupReq    func() *domain.ReportRequest
		wantErr     bool
		errContains string
	}{
		{
			name: "valid request",
			setupReq: func() *domain.ReportRequest {
				return domain.NewReportRequest(userID, domain.ReportTypeDXCC, domain.ReportFormatPDF)
			},
			wantErr: false,
		},
		{
			name: "invalid report type",
			setupReq: func() *domain.ReportRequest {
				return &domain.ReportRequest{
					UserID: userID,
					Type:   domain.ReportType("invalid"),
					Format: domain.ReportFormatPDF,
				}
			},
			wantErr:     true,
			errContains: "invalid report type",
		},
		{
			name: "invalid report format",
			setupReq: func() *domain.ReportRequest {
				return &domain.ReportRequest{
					UserID: userID,
					Type:   domain.ReportTypeDXCC,
					Format: domain.ReportFormat("invalid"),
				}
			},
			wantErr:     true,
			errContains: "invalid report format",
		},
		{
			name: "invalid date range",
			setupReq: func() *domain.ReportRequest {
				req := domain.NewReportRequest(userID, domain.ReportTypeDXCC, domain.ReportFormatPDF)
				req.DateFrom = time.Now()
				req.DateTo = time.Now().Add(-24 * time.Hour)
				return req
			},
			wantErr:     true,
			errContains: "invalid date range",
		},
		{
			name: "valid date range",
			setupReq: func() *domain.ReportRequest {
				req := domain.NewReportRequest(userID, domain.ReportTypeDXCC, domain.ReportFormatPDF)
				req.DateFrom = time.Now().Add(-24 * time.Hour)
				req.DateTo = time.Now()
				return req
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupReq()
			err := req.Validate()

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidReportType(t *testing.T) {
	tests := []struct {
		reportType domain.ReportType
		expected   bool
	}{
		{domain.ReportTypeDXCC, true},
		{domain.ReportTypeAwards, true},
		{domain.ReportTypeContest, true},
		{domain.ReportTypeCustom, true},
		{domain.ReportType("invalid"), false},
		{domain.ReportType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.reportType), func(t *testing.T) {
			result := domain.IsValidReportType(tt.reportType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidReportFormat(t *testing.T) {
	tests := []struct {
		format   domain.ReportFormat
		expected bool
	}{
		{domain.ReportFormatPDF, true},
		{domain.ReportFormatCSV, true},
		{domain.ReportFormatHTML, true},
		{domain.ReportFormatJSON, true},
		{domain.ReportFormat("invalid"), false},
		{domain.ReportFormat(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			result := domain.IsValidReportFormat(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}
