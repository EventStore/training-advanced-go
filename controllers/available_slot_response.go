package controllers

import (
	"time"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
)

type AvailableSlotResponse struct {
	DayId    string
	SlotId   string
	Date     string
	Time     string
	Duration string
}

func AvailableSlotResponseFrom(a readmodel.AvailableSlot) AvailableSlotResponse {
	return AvailableSlotResponse{
		DayId: a.DayId,
		SlotId: a.Id,
		Date: a.Date,
		Time: a.StartTime,
		Duration: time.Time{}.Add(a.Duration).Format("15:04:05"),
	}
}
