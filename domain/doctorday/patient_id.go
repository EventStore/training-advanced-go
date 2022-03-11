package doctorday

type PatientID struct {
	Value string
}

func NewPatientID(id string) PatientID {
	return PatientID{
		Value: id,
	}
}
