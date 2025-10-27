package model

import "github.com/google/uuid"

type ActivityPoint struct {
	Hour  string `json:"hour" db:"hour"`
	Band  string `json:"band" db:"band"`
	Count int    `json:"count" db:"count"`
}

type ActivityHeatmap struct {
	UserID   uuid.UUID       `json:"user_id"`
	Period   string          `json:"period"`
	Activity []ActivityPoint `json:"activity"`
}

func NewActivityHeatmap(userID uuid.UUID, period string) *ActivityHeatmap {
	return &ActivityHeatmap{
		UserID:   userID,
		Period:   period,
		Activity: []ActivityPoint{},
	}
}

func (a *ActivityHeatmap) AddActivityPoint(hour, band string, count int) {
	a.Activity = append(a.Activity, ActivityPoint{
		Hour:  hour,
		Band:  band,
		Count: count,
	})
}

func (a *ActivityHeatmap) GetTotalActivity() int {
	total := 0
	for _, point := range a.Activity {
		total += point.Count
	}
	return total
}

func (a *ActivityHeatmap) GetActivityByBand(band string) int {
	total := 0
	for _, point := range a.Activity {
		if point.Band == band {
			total += point.Count
		}
	}
	return total
}
