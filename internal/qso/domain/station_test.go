package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rd2w/rd2w-log/internal/qso/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStation(t *testing.T) {
	userID := uuid.New()

	validStation := func() *domain.Station {
		station, err := domain.NewStation(
			userID,
			"RD2W",
			"Home Station",
			"Kursk, Russia",
			"KO85",
		)
		require.NoError(t, err)
		return station
	}

	t.Run("valid station", func(t *testing.T) {
		station := validStation()

		assert.NotNil(t, station)
		assert.NotEmpty(t, station.ID)
		assert.Equal(t, userID, station.UserID)
		assert.Equal(t, "RD2W", station.Callsign)
		assert.Equal(t, "Home Station", station.Name)
		assert.Equal(t, "Kursk, Russia", station.Location)
		assert.Equal(t, "KO85", station.GridSquare)
		assert.NotNil(t, station.Equipment)
		assert.Empty(t, station.Equipment)
		assert.NotNil(t, station.Antennas)
		assert.Empty(t, station.Antennas)
		assert.NotNil(t, station.SupportedModes)
		assert.Empty(t, station.SupportedModes)
		assert.NotNil(t, station.SupportedBands)
		assert.Empty(t, station.SupportedBands)
		assert.WithinDuration(t, time.Now(), station.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), station.UpdatedAt, time.Second)
	})

	t.Run("invalid callsign", func(t *testing.T) {
		_, err := domain.NewStation(
			userID,
			"", // empty callsign
			"Home Station",
			"Kursk, Russia",
			"KO85",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid callsign")
	})

	t.Run("invalid name", func(t *testing.T) {
		_, err := domain.NewStation(
			userID,
			"RD2W",
			"", // empty name
			"Kursk, Russia",
			"KO85",
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid station name")
	})

	t.Run("valid with empty location and grid", func(t *testing.T) {
		station, err := domain.NewStation(
			userID,
			"RD2W",
			"Home Station",
			"", // empty location
			"", // empty grid square
		)
		require.NoError(t, err)
		assert.Equal(t, "", station.Location)
		assert.Equal(t, "", station.GridSquare)
	})
}

func TestStation_Validate(t *testing.T) {
	tests := []struct {
		name    string
		station *domain.Station
		wantErr bool
	}{
		{
			name: "valid station",
			station: &domain.Station{
				Callsign: "RD2W",
				Name:     "Home Station",
			},
			wantErr: false,
		},
		{
			name: "empty callsign",
			station: &domain.Station{
				Callsign: "",
				Name:     "Home Station",
			},
			wantErr: true,
		},
		{
			name: "empty name",
			station: &domain.Station{
				Callsign: "RD2W",
				Name:     "",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			station: &domain.Station{
				Callsign: "",
				Name:     "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.station.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStation_AddEquipment(t *testing.T) {
	userID := uuid.New()
	station, err := domain.NewStation(userID, "RD2W", "Home Station", "Location", "KO85")
	require.NoError(t, err)

	oldUpdatedAt := station.UpdatedAt

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	// Add first equipment
	station.AddEquipment("IC-7300")
	assert.Len(t, station.Equipment, 1)
	assert.Contains(t, station.Equipment, "IC-7300")
	assert.True(t, station.UpdatedAt.After(oldUpdatedAt))

	oldUpdatedAt = station.UpdatedAt
	time.Sleep(time.Millisecond * 10)

	// Add second equipment
	station.AddEquipment("LDG Z-100Plus")
	assert.Len(t, station.Equipment, 2)
	assert.Contains(t, station.Equipment, "IC-7300")
	assert.Contains(t, station.Equipment, "LDG Z-100Plus")
	assert.True(t, station.UpdatedAt.After(oldUpdatedAt))
}

func TestStation_AddAntenna(t *testing.T) {
	userID := uuid.New()
	station, err := domain.NewStation(userID, "RD2W", "Home Station", "Location", "KO85")
	require.NoError(t, err)

	oldUpdatedAt := station.UpdatedAt

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	// Add first antenna
	station.AddAntenna("Dipole 40m")
	assert.Len(t, station.Antennas, 1)
	assert.Contains(t, station.Antennas, "Dipole 40m")
	assert.True(t, station.UpdatedAt.After(oldUpdatedAt))

	oldUpdatedAt = station.UpdatedAt
	time.Sleep(time.Millisecond * 10)

	// Add second antenna
	station.AddAntenna("Vertical 20m")
	assert.Len(t, station.Antennas, 2)
	assert.Contains(t, station.Antennas, "Dipole 40m")
	assert.Contains(t, station.Antennas, "Vertical 20m")
	assert.True(t, station.UpdatedAt.After(oldUpdatedAt))
}

func TestStation_AddSupportedMode(t *testing.T) {
	userID := uuid.New()
	station, err := domain.NewStation(userID, "RD2W", "Home Station", "Location", "KO85")
	require.NoError(t, err)

	oldUpdatedAt := station.UpdatedAt

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	// Add first mode
	station.AddSupportedMode("SSB")
	assert.Len(t, station.SupportedModes, 1)
	assert.Contains(t, station.SupportedModes, "SSB")
	assert.True(t, station.UpdatedAt.After(oldUpdatedAt))

	oldUpdatedAt = station.UpdatedAt
	time.Sleep(time.Millisecond * 10)

	// Add more modes
	station.AddSupportedMode("CW")
	station.AddSupportedMode("FT8")
	assert.Len(t, station.SupportedModes, 3)
	assert.Contains(t, station.SupportedModes, "SSB")
	assert.Contains(t, station.SupportedModes, "CW")
	assert.Contains(t, station.SupportedModes, "FT8")
	assert.True(t, station.UpdatedAt.After(oldUpdatedAt))
}

func TestStation_AddSupportedBand(t *testing.T) {
	userID := uuid.New()
	station, err := domain.NewStation(userID, "RD2W", "Home Station", "Location", "KO85")
	require.NoError(t, err)

	oldUpdatedAt := station.UpdatedAt

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	// Add first band
	station.AddSupportedBand("20m")
	assert.Len(t, station.SupportedBands, 1)
	assert.Contains(t, station.SupportedBands, "20m")
	assert.True(t, station.UpdatedAt.After(oldUpdatedAt))

	oldUpdatedAt = station.UpdatedAt
	time.Sleep(time.Millisecond * 10)

	// Add more bands
	station.AddSupportedBand("40m")
	station.AddSupportedBand("80m")
	assert.Len(t, station.SupportedBands, 3)
	assert.Contains(t, station.SupportedBands, "20m")
	assert.Contains(t, station.SupportedBands, "40m")
	assert.Contains(t, station.SupportedBands, "80m")
	assert.True(t, station.UpdatedAt.After(oldUpdatedAt))
}

func TestStation_UpdateDetails(t *testing.T) {
	userID := uuid.New()
	station, err := domain.NewStation(
		userID,
		"RD2W",
		"Old Name",
		"Old Location",
		"KO85",
	)
	require.NoError(t, err)

	oldUpdatedAt := station.UpdatedAt

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	// Update details
	station.UpdateDetails("New Name", "New Location", "KO95")

	assert.Equal(t, "New Name", station.Name)
	assert.Equal(t, "New Location", station.Location)
	assert.Equal(t, "KO95", station.GridSquare)
	assert.True(t, station.UpdatedAt.After(oldUpdatedAt))

	// Verify other fields are unchanged
	assert.Equal(t, "RD2W", station.Callsign)
	assert.Equal(t, userID, station.UserID)
}

func TestStation_MultipleOperations(t *testing.T) {
	userID := uuid.New()
	station, err := domain.NewStation(userID, "RD2W", "Home Station", "Location", "KO85")
	require.NoError(t, err)

	// Perform multiple operations
	station.AddEquipment("IC-7300")
	station.AddAntenna("Dipole 40m")
	station.AddSupportedMode("SSB")
	station.AddSupportedBand("20m")
	station.UpdateDetails("Updated Station", "New Location", "KO95")

	// Verify all operations were applied
	assert.Len(t, station.Equipment, 1)
	assert.Len(t, station.Antennas, 1)
	assert.Len(t, station.SupportedModes, 1)
	assert.Len(t, station.SupportedBands, 1)
	assert.Equal(t, "Updated Station", station.Name)
	assert.Equal(t, "New Location", station.Location)
	assert.Equal(t, "KO95", station.GridSquare)
}

func TestStation_DuplicateAdditions(t *testing.T) {
	userID := uuid.New()
	station, err := domain.NewStation(userID, "RD2W", "Home Station", "Location", "KO85")
	require.NoError(t, err)

	// Add duplicate equipment
	station.AddEquipment("IC-7300")
	station.AddEquipment("IC-7300") // duplicate

	assert.Len(t, station.Equipment, 2) // Should allow duplicates
	assert.Equal(t, "IC-7300", station.Equipment[0])
	assert.Equal(t, "IC-7300", station.Equipment[1])

	// Add duplicate modes
	station.AddSupportedMode("SSB")
	station.AddSupportedMode("SSB") // duplicate

	assert.Len(t, station.SupportedModes, 2) // Should allow duplicates
	assert.Equal(t, "SSB", station.SupportedModes[0])
	assert.Equal(t, "SSB", station.SupportedModes[1])
}

func TestStation_EmptyStringAdditions(t *testing.T) {
	userID := uuid.New()
	station, err := domain.NewStation(userID, "RD2W", "Home Station", "Location", "KO85")
	require.NoError(t, err)

	// Add empty strings (should be allowed)
	station.AddEquipment("")
	station.AddAntenna("")
	station.AddSupportedMode("")
	station.AddSupportedBand("")

	assert.Len(t, station.Equipment, 1)
	assert.Contains(t, station.Equipment, "")
	assert.Len(t, station.Antennas, 1)
	assert.Contains(t, station.Antennas, "")
	assert.Len(t, station.SupportedModes, 1)
	assert.Contains(t, station.SupportedModes, "")
	assert.Len(t, station.SupportedBands, 1)
	assert.Contains(t, station.SupportedBands, "")
}
