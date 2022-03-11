package mongodb

import (
	"time"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
)

type AvailableSlotV2 struct {
	Id        string
	DayId     string
	Date      string
	StartTime string
	Duration  time.Duration
	IsBooked  bool
}

func (a *AvailableSlotV2) ToAvailableSlot() readmodel.AvailableSlot {
	return readmodel.AvailableSlot{
		Id:        a.Id,
		DayId:     a.DayId,
		Date:      a.Date,
		StartTime: a.StartTime,
		Duration:  a.Duration,
	}
}
