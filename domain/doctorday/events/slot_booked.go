package events

import "github.com/google/uuid"

type SlotBooked struct {
	DayId     string    `json:"dayId"`
	SlotId    uuid.UUID `json:"slotId"`
	PatientId string    `json:"patientId"`
}

func NewSlotBooked(slotId uuid.UUID, dayId, patientId string) SlotBooked {
	return SlotBooked{
		DayId:     dayId,
		SlotId:    slotId,
		PatientId: patientId,
	}
}
