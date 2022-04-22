package events


type DayScheduleArchived struct {
	DayId string
}

func NewDayScheduleArchived(dayId string) DayScheduleArchived {
	return DayScheduleArchived{
		DayId: dayId,
	}
}
