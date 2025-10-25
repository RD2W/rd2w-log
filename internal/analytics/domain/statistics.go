package domain

import (
	"time"

	"github.com/google/uuid"
)

type Statistics struct {
	UserID          uuid.UUID      `json:"user_id"`
	TotalQSO        int            `json:"total_qso"`
	QSOByBand       map[string]int `json:"qso_by_band"`
	QSOByMode       map[string]int `json:"qso_by_mode"`
	QSOByContinent  map[string]int `json:"qso_by_continent"`
	Last30Days      []int          `json:"last_30_days"`
	UniqueCountries int            `json:"unique_countries"`
	UniqueGrids     int            `json:"unique_grids"`
	MostActiveBand  string         `json:"most_active_band"`
	MostActiveMode  string         `json:"most_active_mode"`
	GeneratedAt     time.Time      `json:"generated_at"`
}

func NewStatistics(userID uuid.UUID) *Statistics {
	return &Statistics{
		UserID:         userID,
		QSOByBand:      make(map[string]int),
		QSOByMode:      make(map[string]int),
		QSOByContinent: make(map[string]int),
		Last30Days:     make([]int, 30),
		GeneratedAt:    time.Now(),
	}
}

func (s *Statistics) AddQSO(band, mode, continent string) {
	s.TotalQSO++

	// Обновляем статистику по диапазонам
	s.QSOByBand[band]++

	// Обновляем статистику по режимам
	s.QSOByMode[mode]++

	// Обновляем статистику по континентам
	if continent != "" {
		s.QSOByContinent[continent]++
	}

	// Обновляем самые активные диапазон и режим
	s.updateMostActive()
}

func (s *Statistics) AddDailyCount(dayIndex int, count int) {
	if dayIndex >= 0 && dayIndex < len(s.Last30Days) {
		s.Last30Days[dayIndex] = count
	}
}

func (s *Statistics) SetUniqueEntities(countries, grids int) {
	s.UniqueCountries = countries
	s.UniqueGrids = grids
}

func (s *Statistics) updateMostActive() {
	s.MostActiveBand = s.findMaxKey(s.QSOByBand)
	s.MostActiveMode = s.findMaxKey(s.QSOByMode)
}

func (s *Statistics) findMaxKey(m map[string]int) string {
	if len(m) == 0 {
		return ""
	}

	maxKey := ""
	maxValue := -1

	for k, v := range m {
		if v > maxValue {
			maxValue = v
			maxKey = k
		}
	}

	return maxKey
}
