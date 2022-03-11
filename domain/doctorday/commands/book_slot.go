package commands

import "github.com/google/uuid"

type BookSlot struct {
	DayId     string
	SlotId    uuid.UUID
	PatientId string
}

func NewBookSlot(slotId uuid.UUID, dayId, patientId string) BookSlot {
	return BookSlot{
		DayId:     dayId,
		SlotId:    slotId,
		PatientId: patientId,
	}
}
