package commands

type CancelDaySchedule struct {
	DayId string
}

func NewCancelDaySchedule(dayId string) CancelDaySchedule {
	return CancelDaySchedule{
		DayId: dayId,
	}
}
