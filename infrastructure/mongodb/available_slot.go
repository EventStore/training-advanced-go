package mongodb

import (
	"time"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
)

type AvailableSlot struct {
	Id        string
	DayId     string
	Date      string
	StartTime string
	Duration  time.Duration
	IsBooked  bool
}

func (a *AvailableSlot) ToAvailableSlot() readmodel.AvailableSlot {
	return readmodel.AvailableSlot{
		Id:        a.Id,
		DayId:     a.DayId,
		Date:      a.Date,
		StartTime: a.StartTime,
		Duration:  a.Duration,
	}
}
