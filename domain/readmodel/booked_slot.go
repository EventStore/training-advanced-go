package readmodel

type BookedSlot struct {
	SlotId    string
	DayId     string
	Month     int
	PatientId string
	IsBooked  bool
}

func NewBookedSlot(slotId, dayId string, month int) BookedSlot {
	return BookedSlot{
		SlotId: slotId,
		DayId: dayId,
		Month: month,
	}
}
