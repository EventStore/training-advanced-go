package readmodel

import "time"

type AvailableSlot struct {
	Id        string        `json:"id"`
	DayId     string        `json:"dayId"`
	Date      string        `json:"date"`
	StartTime string        `json:"startTime"`
	Duration  time.Duration `json:"duration"`
}

func NewAvailableSlot(id, dayId, date, startTime string, d time.Duration) AvailableSlot {
	return AvailableSlot{
		Id:        id,
		DayId:     dayId,
		Date:      date,
		StartTime: startTime,
		Duration:  d,
	}
}
