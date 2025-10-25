package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReportType string

const (
	ReportTypeDXCC    ReportType = "dxcc"
	ReportTypeAwards  ReportType = "awards"
	ReportTypeContest ReportType = "contest"
	ReportTypeCustom  ReportType = "custom"
)

type ReportFormat string

const (
	ReportFormatPDF  ReportFormat = "pdf"
	ReportFormatCSV  ReportFormat = "csv"
	ReportFormatHTML ReportFormat = "html"
	ReportFormatJSON ReportFormat = "json"
)

type ReportRequest struct {
	UserID   uuid.UUID              `json:"user_id"`
	Type     ReportType             `json:"type"`
	DateFrom time.Time              `json:"date_from"`
	DateTo   time.Time              `json:"date_to"`
	Format   ReportFormat           `json:"format"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type ReportResponse struct {
	Data        []byte       `json:"data"`
	Type        ReportType   `json:"type"`
	Format      ReportFormat `json:"format"`
	Filename    string       `json:"filename"`
	Size        int64        `json:"size"`
	GeneratedAt time.Time    `json:"generated_at"`
}

func NewReportRequest(userID uuid.UUID, reportType ReportType, format ReportFormat) *ReportRequest {
	return &ReportRequest{
		UserID:  userID,
		Type:    reportType,
		Format:  format,
		Options: make(map[string]interface{}),
	}
}

func (r *ReportRequest) Validate() error {
	if !IsValidReportType(r.Type) {
		return ErrInvalidReportType
	}

	if !IsValidReportFormat(r.Format) {
		return ErrInvalidReportFormat
	}

	if !r.DateFrom.IsZero() && !r.DateTo.IsZero() && r.DateFrom.After(r.DateTo) {
		return ErrInvalidDateRange
	}

	return nil
}

func IsValidReportType(t ReportType) bool {
	switch t {
	case ReportTypeDXCC, ReportTypeAwards, ReportTypeContest, ReportTypeCustom:
		return true
	default:
		return false
	}
}

func IsValidReportFormat(f ReportFormat) bool {
	switch f {
	case ReportFormatPDF, ReportFormatCSV, ReportFormatHTML, ReportFormatJSON:
		return true
	default:
		return false
	}
}
