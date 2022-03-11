package commands

type ArchiveDaySchedule struct {
	DayId string
}

func NewArchiveDaySchedule(dayId string) ArchiveDaySchedule {
	return ArchiveDaySchedule{
		DayId: dayId,
	}
}
