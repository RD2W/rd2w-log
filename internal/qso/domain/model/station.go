package model

import (
	"time"

	"github.com/google/uuid"
)

type Station struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Callsign       string    `json:"callsign" db:"callsign"`
	Name           string    `json:"name" db:"name"`
	Location       string    `json:"location" db:"location"`
	GridSquare     string    `json:"grid_square" db:"grid_square"`
	Equipment      []string  `json:"equipment" db:"equipment"`
	Antennas       []string  `json:"antennas" db:"antennas"`
	SupportedModes []string  `json:"supported_modes" db:"supported_modes"`
	SupportedBands []string  `json:"supported_bands" db:"supported_bands"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

func NewStation(
	userID uuid.UUID,
	callsign string,
	name string,
	location string,
	gridSquare string,
) (*Station, error) {

	station := &Station{
		ID:             uuid.New(),
		UserID:         userID,
		Callsign:       callsign,
		Name:           name,
		Location:       location,
		GridSquare:     gridSquare,
		Equipment:      []string{},
		Antennas:       []string{},
		SupportedModes: []string{},
		SupportedBands: []string{},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := station.Validate(); err != nil {
		return nil, err
	}

	return station, nil
}

func (s *Station) Validate() error {
	if s.Callsign == "" {
		return ErrInvalidCallsign
	}

	if s.Name == "" {
		return ErrInvalidStationName
	}

	return nil
}

func (s *Station) AddEquipment(equipment string) {
	s.Equipment = append(s.Equipment, equipment)
	s.UpdatedAt = time.Now()
}

func (s *Station) AddAntenna(antenna string) {
	s.Antennas = append(s.Antennas, antenna)
	s.UpdatedAt = time.Now()
}

func (s *Station) AddSupportedMode(mode string) {
	s.SupportedModes = append(s.SupportedModes, mode)
	s.UpdatedAt = time.Now()
}

func (s *Station) AddSupportedBand(band string) {
	s.SupportedBands = append(s.SupportedBands, band)
	s.UpdatedAt = time.Now()
}

func (s *Station) UpdateDetails(name, location, gridSquare string) {
	s.Name = name
	s.Location = location
	s.GridSquare = gridSquare
	s.UpdatedAt = time.Now()
}
