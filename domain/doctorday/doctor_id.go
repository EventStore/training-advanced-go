package doctorday

import "github.com/google/uuid"

type DoctorID struct {
	Value uuid.UUID
}

func NewDoctorID(id uuid.UUID) DoctorID {
	return DoctorID{
		Value: id,
	}
}
