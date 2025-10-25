package domain

import "errors"

var (
	ErrQSONotFound           = errors.New("qso not found")
	ErrInvalidCallsign       = errors.New("invalid callsign")
	ErrInvalidTime           = errors.New("invalid time")
	ErrInvalidTimeOff        = errors.New("time_off must be after time_on")
	ErrInvalidBand           = errors.New("invalid band")
	ErrInvalidMode           = errors.New("invalid mode")
	ErrInvalidDateRange      = errors.New("invalid date range")
	ErrInvalidPage           = errors.New("invalid page number")
	ErrInvalidLimit          = errors.New("invalid limit")
	ErrStationNotFound       = errors.New("station not found")
	ErrInvalidStationName    = errors.New("invalid station name")
	ErrInvalidQSLSentVia     = errors.New("invalid QSL sent via")
	ErrInvalidQSLReceivedVia = errors.New("invalid QSL received via")
	ErrADIFParseError        = errors.New("failed to parse ADIF file")
	ErrADIFValidationError   = errors.New("ADIF data validation failed")
)
