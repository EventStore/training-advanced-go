package commands

import (
	"time"

	"github.com/google/uuid"
)

type ScheduleDay struct {
	DoctorId uuid.UUID
	Date     time.Time
	Slots    []ScheduledSlot
}

type ScheduledSlot struct {
	StartTime time.Time
	Duration  time.Duration
}

func NewScheduleDay(doctorId uuid.UUID, date time.Time, slots []ScheduledSlot) ScheduleDay {
	return ScheduleDay{
		DoctorId: doctorId,
		Date:     date,
		Slots:    slots,
	}
}

func NewScheduledSlot(start time.Time, duration time.Duration) ScheduledSlot {
	return ScheduledSlot{
		StartTime: start,
		Duration:  duration,
	}
}
