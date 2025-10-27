package model

import (
	"time"

	"github.com/google/uuid"
)

type Filter struct {
	UserID   uuid.UUID
	Band     string
	Mode     string
	Callsign string
	DateFrom time.Time
	DateTo   time.Time
	Page     int
	Limit    int
}

func NewFilter(userID uuid.UUID) *Filter {
	return &Filter{
		UserID: userID,
		Page:   1,
		Limit:  50,
	}
}

func (f *Filter) WithBand(band string) *Filter {
	f.Band = band
	return f
}

func (f *Filter) WithMode(mode string) *Filter {
	f.Mode = mode
	return f
}

func (f *Filter) WithCallsign(callsign string) *Filter {
	f.Callsign = callsign
	return f
}

func (f *Filter) WithDateRange(dateFrom, dateTo time.Time) *Filter {
	f.DateFrom = dateFrom
	f.DateTo = dateTo
	return f
}

func (f *Filter) WithPagination(page, limit int) *Filter {
	f.Page = page
	f.Limit = limit
	return f
}

func (f *Filter) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f *Filter) Validate() error {
	if !f.DateFrom.IsZero() && !f.DateTo.IsZero() && f.DateFrom.After(f.DateTo) {
		return ErrInvalidDateRange
	}

	if f.Page < 1 {
		return ErrInvalidPage
	}

	if f.Limit < 1 || f.Limit > 1000 {
		return ErrInvalidLimit
	}

	return nil
}
