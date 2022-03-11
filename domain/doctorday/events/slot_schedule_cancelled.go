package events

import "github.com/google/uuid"

type SlotScheduleCancelled struct {
	DayId  string    `json:"dayId"`
	SlotId uuid.UUID `json:"slotId"`
}

func NewSlotScheduleCancelled(slotId uuid.UUID, dayId string) SlotScheduleCancelled {
	return SlotScheduleCancelled{
		DayId:  dayId,
		SlotId: slotId,
	}
}
