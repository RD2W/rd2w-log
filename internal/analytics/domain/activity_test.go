package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rd2w/rd2w-log/internal/analytics/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewActivityHeatmap(t *testing.T) {
	userID := uuid.New()
	period := "7d"

	heatmap := domain.NewActivityHeatmap(userID, period)

	assert.NotNil(t, heatmap)
	assert.Equal(t, userID, heatmap.UserID)
	assert.Equal(t, period, heatmap.Period)
	assert.NotNil(t, heatmap.Activity)
	assert.Empty(t, heatmap.Activity)
}

func TestActivityHeatmap_AddActivityPoint(t *testing.T) {
	userID := uuid.New()
	heatmap := domain.NewActivityHeatmap(userID, "7d")

	// Add first activity point
	heatmap.AddActivityPoint("14:00", "20m", 5)
	assert.Len(t, heatmap.Activity, 1)
	assert.Equal(t, "14:00", heatmap.Activity[0].Hour)
	assert.Equal(t, "20m", heatmap.Activity[0].Band)
	assert.Equal(t, 5, heatmap.Activity[0].Count)

	// Add second activity point
	heatmap.AddActivityPoint("15:00", "40m", 3)
	assert.Len(t, heatmap.Activity, 2)
	assert.Equal(t, "15:00", heatmap.Activity[1].Hour)
	assert.Equal(t, "40m", heatmap.Activity[1].Band)
	assert.Equal(t, 3, heatmap.Activity[1].Count)

	// Add point with same hour and band (should be allowed)
	heatmap.AddActivityPoint("14:00", "20m", 2)
	assert.Len(t, heatmap.Activity, 3)
	assert.Equal(t, "14:00", heatmap.Activity[2].Hour)
	assert.Equal(t, "20m", heatmap.Activity[2].Band)
	assert.Equal(t, 2, heatmap.Activity[2].Count)
}

func TestActivityHeatmap_GetTotalActivity(t *testing.T) {
	userID := uuid.New()
	heatmap := domain.NewActivityHeatmap(userID, "7d")

	// Test empty heatmap
	total := heatmap.GetTotalActivity()
	assert.Equal(t, 0, total)

	// Test with single point
	heatmap.AddActivityPoint("14:00", "20m", 5)
	total = heatmap.GetTotalActivity()
	assert.Equal(t, 5, total)

	// Test with multiple points
	heatmap.AddActivityPoint("15:00", "40m", 3)
	heatmap.AddActivityPoint("16:00", "80m", 7)
	total = heatmap.GetTotalActivity()
	assert.Equal(t, 15, total) // 5 + 3 + 7

	// Test with zero counts
	heatmap.AddActivityPoint("17:00", "20m", 0)
	total = heatmap.GetTotalActivity()
	assert.Equal(t, 15, total) // 0 should not affect total

	// Test with negative counts (if supported)
	heatmap.AddActivityPoint("18:00", "40m", -2)
	total = heatmap.GetTotalActivity()
	assert.Equal(t, 13, total) // 15 - 2
}

func TestActivityHeatmap_GetActivityByBand(t *testing.T) {
	userID := uuid.New()
	heatmap := domain.NewActivityHeatmap(userID, "7d")

	// Setup test data
	heatmap.AddActivityPoint("14:00", "20m", 5)
	heatmap.AddActivityPoint("15:00", "40m", 3)
	heatmap.AddActivityPoint("16:00", "20m", 2)
	heatmap.AddActivityPoint("17:00", "80m", 7)
	heatmap.AddActivityPoint("18:00", "40m", 4)

	tests := []struct {
		name     string
		band     string
		expected int
	}{
		{
			name:     "20m band",
			band:     "20m",
			expected: 7, // 5 + 2
		},
		{
			name:     "40m band",
			band:     "40m",
			expected: 7, // 3 + 4
		},
		{
			name:     "80m band",
			band:     "80m",
			expected: 7, // 7
		},
		{
			name:     "non-existent band",
			band:     "160m",
			expected: 0,
		},
		{
			name:     "empty band",
			band:     "",
			expected: 0,
		},
		{
			name:     "case sensitive band",
			band:     "20M", // different case
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := heatmap.GetActivityByBand(tt.band)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestActivityHeatmap_MixedScenarios(t *testing.T) {
	userID := uuid.New()
	heatmap := domain.NewActivityHeatmap(userID, "30d")

	// Test complex scenario with various data
	testData := []struct {
		hour  string
		band  string
		count int
	}{
		{"00:00", "80m", 2},
		{"01:00", "40m", 1},
		{"02:00", "20m", 3},
		{"12:00", "20m", 8},
		{"13:00", "15m", 5},
		{"14:00", "10m", 4},
		{"15:00", "20m", 6},
		{"16:00", "40m", 2},
		{"17:00", "80m", 1},
		{"18:00", "40m", 3},
		{"23:00", "160m", 1},
	}

	// Add all test data
	for _, data := range testData {
		heatmap.AddActivityPoint(data.hour, data.band, data.count)
	}

	// Verify total activity
	total := heatmap.GetTotalActivity()
	expectedTotal := 2 + 1 + 3 + 8 + 5 + 4 + 6 + 2 + 1 + 3 + 1
	assert.Equal(t, expectedTotal, total)

	// Verify band-specific activities
	assert.Equal(t, 3+8+6, heatmap.GetActivityByBand("20m"))      // 17
	assert.Equal(t, 1+2+3, heatmap.GetActivityByBand("40m"))      // 6
	assert.Equal(t, 2+1, heatmap.GetActivityByBand("80m"))        // 3
	assert.Equal(t, 5, heatmap.GetActivityByBand("15m"))          // 5
	assert.Equal(t, 4, heatmap.GetActivityByBand("10m"))          // 4
	assert.Equal(t, 1, heatmap.GetActivityByBand("160m"))         // 1
	assert.Equal(t, 0, heatmap.GetActivityByBand("non-existent")) // 0
}

func TestActivityHeatmap_EdgeCases(t *testing.T) {
	userID := uuid.New()
	heatmap := domain.NewActivityHeatmap(userID, "1d")

	// Test empty strings
	heatmap.AddActivityPoint("", "", 5)
	assert.Equal(t, "", heatmap.Activity[0].Hour)
	assert.Equal(t, "", heatmap.Activity[0].Band)
	assert.Equal(t, 5, heatmap.Activity[0].Count)

	// Test very large numbers
	heatmap.AddActivityPoint("12:00", "20m", 1000000)
	assert.Equal(t, 1000000, heatmap.Activity[1].Count)

	// Test negative numbers
	heatmap.AddActivityPoint("13:00", "40m", -10)
	assert.Equal(t, -10, heatmap.Activity[2].Count)

	// Verify totals with edge cases
	total := heatmap.GetTotalActivity()
	assert.Equal(t, 5+1000000-10, total) // 999995

	// Test band activity with empty band
	emptyBandActivity := heatmap.GetActivityByBand("")
	assert.Equal(t, 5, emptyBandActivity)
}

func TestActivityPoint_Structure(t *testing.T) {
	// Test that ActivityPoint fields are properly exported and accessible
	point := domain.ActivityPoint{
		Hour:  "14:00",
		Band:  "20m",
		Count: 5,
	}

	assert.Equal(t, "14:00", point.Hour)
	assert.Equal(t, "20m", point.Band)
	assert.Equal(t, 5, point.Count)
}
