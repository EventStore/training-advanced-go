package doctorday

import "github.com/google/uuid"

type SlotID struct {
	Value uuid.UUID
}

func NewSlotID(id uuid.UUID) SlotID {
	return SlotID{
		Value: id,
	}
}
