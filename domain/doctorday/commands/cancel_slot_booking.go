package commands

import (
	"github.com/google/uuid"
)

type CancelSlotBooking struct {
	DayId  string
	SlotId uuid.UUID
	Reason string
}

func NewCancelSlotBooking(slotId uuid.UUID, dayId, reason string) CancelSlotBooking {
	return CancelSlotBooking{
		DayId:  dayId,
		SlotId: slotId,
		Reason: reason,
	}
}
