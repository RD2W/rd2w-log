package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rd2w/rd2w-log/internal/qso/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFilter(t *testing.T) {
	userID := uuid.New()
	filter := domain.NewFilter(userID)

	assert.Equal(t, userID, filter.UserID)
	assert.Equal(t, 1, filter.Page)
	assert.Equal(t, 50, filter.Limit)
	assert.Empty(t, filter.Band)
	assert.Empty(t, filter.Mode)
	assert.Empty(t, filter.Callsign)
	assert.True(t, filter.DateFrom.IsZero())
	assert.True(t, filter.DateTo.IsZero())
}

func TestFilter_WithMethods(t *testing.T) {
	userID := uuid.New()
	dateFrom := time.Now().AddDate(0, -1, 0)
	dateTo := time.Now()

	filter := domain.NewFilter(userID).
		WithBand("20m").
		WithMode("SSB").
		WithCallsign("UA3ABC").
		WithDateRange(dateFrom, dateTo).
		WithPagination(2, 100)

	assert.Equal(t, "20m", filter.Band)
	assert.Equal(t, "SSB", filter.Mode)
	assert.Equal(t, "UA3ABC", filter.Callsign)
	assert.Equal(t, dateFrom, filter.DateFrom)
	assert.Equal(t, dateTo, filter.DateTo)
	assert.Equal(t, 2, filter.Page)
	assert.Equal(t, 100, filter.Limit)
}

func TestFilter_Offset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expected int
	}{
		{"page 1, limit 50", 1, 50, 0},
		{"page 2, limit 50", 2, 50, 50},
		{"page 3, limit 20", 3, 20, 40},
		{"page 1, limit 100", 1, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := uuid.New()
			filter := domain.NewFilter(userID).WithPagination(tt.page, tt.limit)
			assert.Equal(t, tt.expected, filter.Offset())
		})
	}
}

func TestFilter_Validate(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		setupFilter func() *domain.Filter
		wantErr     bool
		errContains string
	}{
		{
			name: "valid filter",
			setupFilter: func() *domain.Filter {
				return domain.NewFilter(userID)
			},
			wantErr: false,
		},
		{
			name: "valid filter with date range",
			setupFilter: func() *domain.Filter {
				return domain.NewFilter(userID).
					WithDateRange(
						time.Now().AddDate(0, -1, 0),
						time.Now(),
					)
			},
			wantErr: false,
		},
		{
			name: "invalid date range",
			setupFilter: func() *domain.Filter {
				return domain.NewFilter(userID).
					WithDateRange(
						time.Now(),
						time.Now().AddDate(0, -1, 0),
					)
			},
			wantErr:     true,
			errContains: "invalid date range",
		},
		{
			name: "invalid page - zero",
			setupFilter: func() *domain.Filter {
				return domain.NewFilter(userID).WithPagination(0, 50)
			},
			wantErr:     true,
			errContains: "invalid page",
		},
		{
			name: "invalid page - negative",
			setupFilter: func() *domain.Filter {
				return domain.NewFilter(userID).WithPagination(-1, 50)
			},
			wantErr:     true,
			errContains: "invalid page",
		},
		{
			name: "invalid limit - zero",
			setupFilter: func() *domain.Filter {
				return domain.NewFilter(userID).WithPagination(1, 0)
			},
			wantErr:     true,
			errContains: "invalid limit",
		},
		{
			name: "invalid limit - negative",
			setupFilter: func() *domain.Filter {
				return domain.NewFilter(userID).WithPagination(1, -5)
			},
			wantErr:     true,
			errContains: "invalid limit",
		},
		{
			name: "invalid limit - too large",
			setupFilter: func() *domain.Filter {
				return domain.NewFilter(userID).WithPagination(1, 1001)
			},
			wantErr:     true,
			errContains: "invalid limit",
		},
		{
			name: "valid limit - boundary values",
			setupFilter: func() *domain.Filter {
				return domain.NewFilter(userID).WithPagination(1, 1) // min valid
			},
			wantErr: false,
		},
		{
			name: "valid limit - max boundary",
			setupFilter: func() *domain.Filter {
				return domain.NewFilter(userID).WithPagination(1, 1000) // max valid
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := tt.setupFilter()
			err := filter.Validate()

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
