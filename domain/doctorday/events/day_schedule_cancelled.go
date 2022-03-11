package events

type DayScheduleCancelled struct {
	DayId string
}

func NewDayScheduleCancelled(dayId string) DayScheduleCancelled {
	return DayScheduleCancelled{
		DayId: dayId,
	}
}
