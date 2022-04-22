package events

import "github.com/google/uuid"

type SlotBookingCancelled struct {
	DayId  string    `json:"dayId"`
	SlotId uuid.UUID `json:"slotId"`
	Reason string    `json:"reason"`
}

func NewSlotBookingCancelled(slotId uuid.UUID, dayId, reason string) SlotBookingCancelled {
	return SlotBookingCancelled{
		DayId:  dayId,
		SlotId: slotId,
		Reason: reason,
	}
}
