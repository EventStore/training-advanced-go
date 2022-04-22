package commands

import (
	"time"

	"github.com/google/uuid"
)

type ScheduleSlot struct {
	SlotId   uuid.UUID
	DoctorId uuid.UUID
	Start    time.Time
	Duration time.Duration
}

func NewScheduleSlot(slotId, doctorId uuid.UUID, start time.Time, duration time.Duration) ScheduleSlot {
	return ScheduleSlot{
		SlotId:   slotId,
		DoctorId: doctorId,
		Start:    start,
		Duration: duration,
	}
}
