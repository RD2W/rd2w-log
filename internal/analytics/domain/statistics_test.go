package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rd2w/rd2w-log/internal/analytics/domain"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
)

func TestNewStatistics(t *testing.T) {
	userID := uuid.New()
	stats := domain.NewStatistics(userID)

	assert.Equal(t, userID, stats.UserID)
	assert.Equal(t, 0, stats.TotalQSO)
	assert.NotNil(t, stats.QSOByBand)
	assert.NotNil(t, stats.QSOByMode)
	assert.NotNil(t, stats.QSOByContinent)
	assert.Len(t, stats.Last30Days, 30)
	assert.Equal(t, 0, stats.UniqueCountries)
	assert.Equal(t, 0, stats.UniqueGrids)
	assert.Empty(t, stats.MostActiveBand)
	assert.Empty(t, stats.MostActiveMode)
}

func TestStatistics_AddQSO(t *testing.T) {
	userID := uuid.New()
	stats := domain.NewStatistics(userID)

	// Add first QSO
	stats.AddQSO("20m", "SSB", "EU")
	assert.Equal(t, 1, stats.TotalQSO)
	assert.Equal(t, 1, stats.QSOByBand["20m"])
	assert.Equal(t, 1, stats.QSOByMode["SSB"])
	assert.Equal(t, 1, stats.QSOByContinent["EU"])
	assert.Equal(t, "20m", stats.MostActiveBand)
	assert.Equal(t, "SSB", stats.MostActiveMode)

	// Add second QSO on same band/mode
	stats.AddQSO("20m", "SSB", "EU")
	assert.Equal(t, 2, stats.TotalQSO)
	assert.Equal(t, 2, stats.QSOByBand["20m"])
	assert.Equal(t, 2, stats.QSOByMode["SSB"])
	assert.Equal(t, 2, stats.QSOByContinent["EU"])

	// Add QSO on different band/mode
	stats.AddQSO("40m", "CW", "NA")
	assert.Equal(t, 3, stats.TotalQSO)
	assert.Equal(t, 1, stats.QSOByBand["40m"])
	assert.Equal(t, 1, stats.QSOByMode["CW"])
	assert.Equal(t, 1, stats.QSOByContinent["NA"])

	// Add QSO without continent
	stats.AddQSO("20m", "FT8", "")
	assert.Equal(t, 4, stats.TotalQSO)
	assert.Equal(t, 3, stats.QSOByBand["20m"]) // 20m should still be most active
	assert.Equal(t, 1, stats.QSOByMode["FT8"])
}

func TestStatistics_AddDailyCount(t *testing.T) {
	userID := uuid.New()
	stats := domain.NewStatistics(userID)

	// Valid indices
	stats.AddDailyCount(0, 5)
	stats.AddDailyCount(29, 10)
	assert.Equal(t, 5, stats.Last30Days[0])
	assert.Equal(t, 10, stats.Last30Days[29])

	// Invalid indices should be ignored
	stats.AddDailyCount(-1, 100)
	stats.AddDailyCount(30, 100)
	// No change should occur for invalid indices
}

func TestStatistics_SetUniqueEntities(t *testing.T) {
	userID := uuid.New()
	stats := domain.NewStatistics(userID)

	stats.SetUniqueEntities(25, 50)
	assert.Equal(t, 25, stats.UniqueCountries)
	assert.Equal(t, 50, stats.UniqueGrids)
}

func TestStatistics_UpdateMostActive(t *testing.T) {
	userID := uuid.New()

	// Test empty statistics
	stats := domain.NewStatistics(userID)
	stats.AddQSO("", "", "") // This will trigger updateMostActive internally
	assert.Equal(t, "", stats.MostActiveBand)
	assert.Equal(t, "", stats.MostActiveMode)

	// Test single band/mode
	stats = domain.NewStatistics(userID)
	stats.AddQSO("20m", "SSB", "EU")
	assert.Equal(t, "20m", stats.MostActiveBand)
	assert.Equal(t, "SSB", stats.MostActiveMode)

	// Test multiple bands - 20m should remain most active
	stats = domain.NewStatistics(userID)
	stats.AddQSO("20m", "SSB", "EU")
	stats.AddQSO("20m", "SSB", "EU")
	stats.AddQSO("40m", "CW", "NA")
	assert.Equal(t, "20m", stats.MostActiveBand) // 20m has 2, 40m has 1
	assert.Equal(t, "SSB", stats.MostActiveMode) // SSB has 2, CW has 1

	// Test when new band becomes most active
	stats = domain.NewStatistics(userID)
	stats.AddQSO("20m", "SSB", "EU")
	stats.AddQSO("40m", "CW", "NA")
	stats.AddQSO("40m", "CW", "NA")
	stats.AddQSO("40m", "CW", "NA")
	// Now 20m: 1, 40m: 3 - 40m should become most active
	assert.Equal(t, "40m", stats.MostActiveBand)
	assert.Equal(t, "CW", stats.MostActiveMode)
}

func TestStatistics_FindMaxKeyBehavior(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name         string
		addQSOs      []struct{ band, mode, continent string }
		expectedBand string
		expectedMode string
	}{
		{
			name:         "empty",
			addQSOs:      []struct{ band, mode, continent string }{},
			expectedBand: "",
			expectedMode: "",
		},
		{
			name: "single element",
			addQSOs: []struct{ band, mode, continent string }{
				{"20m", "SSB", "EU"},
			},
			expectedBand: "20m",
			expectedMode: "SSB",
		},
		{
			name: "clear winner band",
			addQSOs: []struct{ band, mode, continent string }{
				{"20m", "SSB", "EU"},
				{"40m", "CW", "NA"},
				{"40m", "CW", "NA"}, // 40m has 2, 20m has 1
			},
			expectedBand: "40m",
			expectedMode: "CW",
		},
		{
			name: "clear winner mode",
			addQSOs: []struct{ band, mode, continent string }{
				{"20m", "SSB", "EU"},
				{"20m", "FT8", "EU"},
				{"20m", "FT8", "EU"}, // FT8 has 2, SSB has 1
			},
			expectedBand: "20m",
			expectedMode: "FT8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := domain.NewStatistics(userID)

			for _, qso := range tt.addQSOs {
				stats.AddQSO(qso.band, qso.mode, qso.continent)
			}

			assert.Equal(t, tt.expectedBand, stats.MostActiveBand)
			assert.Equal(t, tt.expectedMode, stats.MostActiveMode)
		})
	}
}

func TestStatistics_TieBehavior(t *testing.T) {
	userID := uuid.New()

	// When there's a tie, the result is non-deterministic due to random map iteration in Go.
	// However, we can verify that one of the leaders is selected and counters are correct.
	stats := domain.NewStatistics(userID)

	// Create a tie situation - both bands and modes have 2 QSOs each
	stats.AddQSO("20m", "SSB", "EU")
	stats.AddQSO("20m", "SSB", "EU")
	stats.AddQSO("40m", "CW", "NA")
	stats.AddQSO("40m", "CW", "NA")

	// Verify that one of the leaders is selected (but we don't care which specific one)
	// This accommodates Go's random map iteration behavior while ensuring business logic works
	assert.True(t, stats.MostActiveBand == "20m" || stats.MostActiveBand == "40m")
	assert.True(t, stats.MostActiveMode == "SSB" || stats.MostActiveMode == "CW")

	// Verify that the counters are correct regardless of which leader was selected
	assert.Equal(t, 2, stats.QSOByBand["20m"])
	assert.Equal(t, 2, stats.QSOByBand["40m"])
	assert.Equal(t, 2, stats.QSOByMode["SSB"])
	assert.Equal(t, 2, stats.QSOByMode["CW"])
}

func TestStatistics_GeneratedAt(t *testing.T) {
	userID := uuid.New()
	before := time.Now()
	stats := domain.NewStatistics(userID)
	after := time.Now()

	assert.True(t, stats.GeneratedAt.After(before) || stats.GeneratedAt.Equal(before))
	assert.True(t, stats.GeneratedAt.Before(after) || stats.GeneratedAt.Equal(after))
}
