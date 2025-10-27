package model_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rd2w/rd2w-log/internal/analytics/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReportRequest(t *testing.T) {
	userID := uuid.New()
	reportType := model.ReportTypeDXCC
	format := model.ReportFormatPDF

	req := model.NewReportRequest(userID, reportType, format)

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
		setupReq    func() *model.ReportRequest
		wantErr     bool
		errContains string
	}{
		{
			name: "valid request",
			setupReq: func() *model.ReportRequest {
				return model.NewReportRequest(userID, model.ReportTypeDXCC, model.ReportFormatPDF)
			},
			wantErr: false,
		},
		{
			name: "invalid report type",
			setupReq: func() *model.ReportRequest {
				return &model.ReportRequest{
					UserID: userID,
					Type:   model.ReportType("invalid"),
					Format: model.ReportFormatPDF,
				}
			},
			wantErr:     true,
			errContains: "invalid report type",
		},
		{
			name: "invalid report format",
			setupReq: func() *model.ReportRequest {
				return &model.ReportRequest{
					UserID: userID,
					Type:   model.ReportTypeDXCC,
					Format: model.ReportFormat("invalid"),
				}
			},
			wantErr:     true,
			errContains: "invalid report format",
		},
		{
			name: "invalid date range",
			setupReq: func() *model.ReportRequest {
				req := model.NewReportRequest(userID, model.ReportTypeDXCC, model.ReportFormatPDF)
				req.DateFrom = time.Now()
				req.DateTo = time.Now().Add(-24 * time.Hour)
				return req
			},
			wantErr:     true,
			errContains: "invalid date range",
		},
		{
			name: "valid date range",
			setupReq: func() *model.ReportRequest {
				req := model.NewReportRequest(userID, model.ReportTypeDXCC, model.ReportFormatPDF)
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
		reportType model.ReportType
		expected   bool
	}{
		{model.ReportTypeDXCC, true},
		{model.ReportTypeAwards, true},
		{model.ReportTypeContest, true},
		{model.ReportTypeCustom, true},
		{model.ReportType("invalid"), false},
		{model.ReportType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.reportType), func(t *testing.T) {
			result := model.IsValidReportType(tt.reportType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidReportFormat(t *testing.T) {
	tests := []struct {
		format   model.ReportFormat
		expected bool
	}{
		{model.ReportFormatPDF, true},
		{model.ReportFormatCSV, true},
		{model.ReportFormatHTML, true},
		{model.ReportFormatJSON, true},
		{model.ReportFormat("invalid"), false},
		{model.ReportFormat(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			result := model.IsValidReportFormat(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}
