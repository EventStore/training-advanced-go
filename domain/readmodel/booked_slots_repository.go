package readmodel

type BookedSlotsRepository interface {
	AddSlot(s BookedSlot) error
	CountByPatientAndMonth(patientId string, month int) (int, error)
	MarkSlotAsAvailable(slotId string) error
	MarkSlotAsBooked(slotId, patientId string) error
	GetSlot(slotId string) (BookedSlot, error)
}
