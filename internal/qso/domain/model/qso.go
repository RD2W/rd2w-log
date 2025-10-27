package model

import (
	"time"

	"github.com/google/uuid"
)

type QSO struct {
	ID              uuid.UUID `json:"id" db:"id"`
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	Callsign        string    `json:"callsign" db:"callsign"`
	TimeOn          time.Time `json:"time_on" db:"time_on"`
	TimeOff         time.Time `json:"time_off" db:"time_off"`
	Frequency       string    `json:"frequency" db:"frequency"`
	Band            string    `json:"band" db:"band"`
	Mode            string    `json:"mode" db:"mode"`
	RSTSent         string    `json:"rst_sent" db:"rst_sent"`
	RSTReceived     string    `json:"rst_received" db:"rst_received"`
	GridSquare      string    `json:"grid_square" db:"grid_square"`
	Operator        string    `json:"operator" db:"operator"`
	StationCallsign string    `json:"station_callsign" db:"station_callsign"`
	Contest         string    `json:"contest" db:"contest"`
	Notes           string    `json:"notes" db:"notes"`
	QSL             *QSLInfo  `json:"qsl" db:"-"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func NewQSO(
	userID uuid.UUID,
	callsign string,
	timeOn time.Time,
	frequency string,
	band string,
	mode string,
	rstSent string,
	rstReceived string,
) (*QSO, error) {

	qso := &QSO{
		ID:          uuid.New(),
		UserID:      userID,
		Callsign:    callsign,
		TimeOn:      timeOn,
		TimeOff:     timeOn,
		Frequency:   frequency,
		Band:        band,
		Mode:        mode,
		RSTSent:     rstSent,
		RSTReceived: rstReceived,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := qso.Validate(); err != nil {
		return nil, err
	}

	return qso, nil
}

func (q *QSO) Validate() error {
	if q.Callsign == "" {
		return ErrInvalidCallsign
	}

	if q.TimeOn.IsZero() {
		return ErrInvalidTime
	}

	if q.Band == "" {
		return ErrInvalidBand
	}

	if q.Mode == "" {
		return ErrInvalidMode
	}

	if q.TimeOff.Before(q.TimeOn) {
		return ErrInvalidTimeOff
	}

	return nil
}

func (q *QSO) UpdateDetails(
	rstSent string,
	rstReceived string,
	gridSquare string,
	operator string,
	stationCallsign string,
	contest string,
	notes string,
) {
	q.RSTSent = rstSent
	q.RSTReceived = rstReceived
	q.GridSquare = gridSquare
	q.Operator = operator
	q.StationCallsign = stationCallsign
	q.Contest = contest
	q.Notes = notes
	q.UpdatedAt = time.Now()
}

func (q *QSO) SetTimeOff(timeOff time.Time) error {
	if timeOff.Before(q.TimeOn) {
		return ErrInvalidTimeOff
	}

	q.TimeOff = timeOff
	q.UpdatedAt = time.Now()
	return nil
}

func (q *QSO) SetQSLInfo(qslInfo *QSLInfo) {
	q.QSL = qslInfo
}
