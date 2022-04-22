package controllers

import (
	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/google/uuid"
)

type CancelSlotBookingRequest struct {
	SlotId uuid.UUID `json:"slotId"`
	Reason string    `json:"reason"`
}

func (r *CancelSlotBookingRequest) ToCommand(dayId string) commands.CancelSlotBooking {
	return commands.NewCancelSlotBooking(r.SlotId, dayId, r.Reason)
}
