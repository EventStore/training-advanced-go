package controllers

import (
	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/google/uuid"
)

type BookSlotRequest struct {
	SlotId    uuid.UUID `json:"slotId"`
	PatientId string    `json:"patientId"`
}

func (r *BookSlotRequest) ToCommand(dayId string) commands.BookSlot {
	return commands.NewBookSlot(r.SlotId, dayId, r.PatientId)
}
