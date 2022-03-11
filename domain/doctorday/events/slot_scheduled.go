package events

import (
	"time"

	"github.com/google/uuid"
)

type SlotScheduled struct {
	SlotId    uuid.UUID     `json:"slotId"`
	DayId     string        `json:"dayId"`
	StartTime time.Time     `json:"startTime"`
	Duration  time.Duration `json:"duration"`
}

func NewSlotScheduled(slotId uuid.UUID, dayId string, start time.Time, duration time.Duration) SlotScheduled {
	return SlotScheduled{
		SlotId:    slotId,
		DayId:     dayId,
		StartTime: start,
		Duration:  duration,
	}
}
