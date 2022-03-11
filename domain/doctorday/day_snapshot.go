package doctorday

import (
	"time"

	"github.com/google/uuid"
)

type DaySnapshot struct {
	IsArchived  bool `json:"isArchived"`
	IsCancelled bool `json:"isCancelled"`
	IsScheduled bool `json:"isScheduled"`

	Slots []SlotSnapshot `json:"slots"`
}

func NewDaySnapshot(isArchived, isCancelled, isScheduled bool, slots []SlotSnapshot) DaySnapshot {
	return DaySnapshot{
		IsArchived:  isArchived,
		IsCancelled: isCancelled,
		IsScheduled: isScheduled,
		Slots:       slots,
	}
}

type SlotSnapshot struct {
	Id       uuid.UUID     `json:"id,omitempty"`
	Start    time.Time     `json:"start"`
	Duration time.Duration `json:"duration,omitempty"`
	Booked   bool          `json:"booked,omitempty"`
}

func NewSlotSnapshot(id uuid.UUID, start time.Time, duration time.Duration, booked bool) SlotSnapshot {
	return SlotSnapshot{
		Id:       id,
		Start:    start,
		Duration: duration,
		Booked:   booked,
	}
}
