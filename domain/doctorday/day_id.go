package doctorday

import "time"

type DayID struct {
	Value string
}

func NewDayID(doctorId DoctorID, date time.Time) DayID {
	return DayID{
		Value: doctorId.Value.String() + "_" + date.Format("2006-01-02"),
	}
}

func NewDayIDFrom(dayId string) DayID {
	return DayID{
		Value: dayId,
	}
}
