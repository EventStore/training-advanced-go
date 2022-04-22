package readmodel

import "time"

type ScheduledSlot struct {
	ScheduledSlotId string
	StartTime       time.Time
	Duration        time.Duration
}

func NewScheduledSlot(scheduledSlotId string, startTime time.Time, duration time.Duration) ScheduledSlot {
	return ScheduledSlot{
		ScheduledSlotId: scheduledSlotId,
		StartTime: startTime,
		Duration: duration,
	}
}
