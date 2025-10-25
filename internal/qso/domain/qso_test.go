package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rd2w/rd2w-log/internal/qso/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQSO(t *testing.T) {
	userID := uuid.New()
	timeOn := time.Now()

	validQSO := func() *domain.QSO {
		qso, err := domain.NewQSO(
			userID,
			"UA3ABC",
			timeOn,
			"14.250",
			"20m",
			"SSB",
			"59",
			"59",
		)
		require.NoError(t, err)
		return qso
	}

	t.Run("valid QSO", func(t *testing.T) {
		qso := validQSO()

		assert.NotNil(t, qso)
		assert.NotEmpty(t, qso.ID)
		assert.Equal(t, userID, qso.UserID)
		assert.Equal(t, "UA3ABC", qso.Callsign)
		assert.Equal(t, timeOn, qso.TimeOn)
		assert.Equal(t, timeOn, qso.TimeOff) // TimeOff should equal TimeOn by default
		assert.Equal(t, "14.250", qso.Frequency)
		assert.Equal(t, "20m", qso.Band)
		assert.Equal(t, "SSB", qso.Mode)
		assert.Equal(t, "59", qso.RSTSent)
		assert.Equal(t, "59", qso.RSTReceived)
		assert.WithinDuration(t, time.Now(), qso.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), qso.UpdatedAt, time.Second)
	})

	t.Run("invalid callsign", func(t *testing.T) {
		_, err := domain.NewQSO(
			userID,
			"",
			timeOn,
			"14.250",
			"20m",
			"SSB",
			"59",
			"59",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid callsign")
	})

	t.Run("invalid time", func(t *testing.T) {
		_, err := domain.NewQSO(
			userID,
			"UA3ABC",
			time.Time{},
			"14.250",
			"20m",
			"SSB",
			"59",
			"59",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid time")
	})

	t.Run("invalid band", func(t *testing.T) {
		_, err := domain.NewQSO(
			userID,
			"UA3ABC",
			timeOn,
			"14.250",
			"",
			"SSB",
			"59",
			"59",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid band")
	})

	t.Run("invalid mode", func(t *testing.T) {
		_, err := domain.NewQSO(
			userID,
			"UA3ABC",
			timeOn,
			"14.250",
			"20m",
			"",
			"59",
			"59",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid mode")
	})
}

func TestQSO_UpdateDetails(t *testing.T) {
	userID := uuid.New()
	qso, err := domain.NewQSO(
		userID,
		"UA3ABC",
		time.Now(),
		"14.250",
		"20m",
		"SSB",
		"59",
		"59",
	)
	require.NoError(t, err)

	oldUpdatedAt := qso.UpdatedAt

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	qso.UpdateDetails(
		"57",
		"55",
		"KO85",
		"RD2W",
		"RD2W",
		"Contest",
		"Good conditions",
	)

	assert.Equal(t, "57", qso.RSTSent)
	assert.Equal(t, "55", qso.RSTReceived)
	assert.Equal(t, "KO85", qso.GridSquare)
	assert.Equal(t, "RD2W", qso.Operator)
	assert.Equal(t, "RD2W", qso.StationCallsign)
	assert.Equal(t, "Contest", qso.Contest)
	assert.Equal(t, "Good conditions", qso.Notes)
	assert.True(t, qso.UpdatedAt.After(oldUpdatedAt))
}

func TestQSO_SetTimeOff(t *testing.T) {
	userID := uuid.New()
	timeOn := time.Now()
	qso, err := domain.NewQSO(
		userID,
		"UA3ABC",
		timeOn,
		"14.250",
		"20m",
		"SSB",
		"59",
		"59",
	)
	require.NoError(t, err)

	oldUpdatedAt := qso.UpdatedAt
	timeOff := timeOn.Add(5 * time.Minute)

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	err = qso.SetTimeOff(timeOff)
	require.NoError(t, err)

	assert.Equal(t, timeOff, qso.TimeOff)
	assert.True(t, qso.UpdatedAt.After(oldUpdatedAt))

	// Test invalid time (timeOff before timeOn)
	err = qso.SetTimeOff(timeOn.Add(-time.Minute))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "time_off must be after time_on")
}

func TestQSO_SetQSLInfo(t *testing.T) {
	userID := uuid.New()
	qso, err := domain.NewQSO(
		userID,
		"UA3ABC",
		time.Now(),
		"14.250",
		"20m",
		"SSB",
		"59",
		"59",
	)
	require.NoError(t, err)

	assert.Nil(t, qso.QSL)

	qslInfo := domain.NewQSLInfo(qso.ID)
	qso.SetQSLInfo(qslInfo)

	assert.NotNil(t, qso.QSL)
	assert.Equal(t, qso.ID, qso.QSL.QSOID)
}
