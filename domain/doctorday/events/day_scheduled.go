package events

import (
	"time"

	"github.com/google/uuid"
)

type DayScheduled struct {
	DayId    string    `json:"dayId"`
	DoctorId uuid.UUID `json:"doctorId"`
	Date     time.Time `json:"date"`
}

func NewDayScheduled(dayId string, doctorId uuid.UUID, date time.Time) DayScheduled {
	return DayScheduled{
		DayId:    dayId,
		DoctorId: doctorId,
		Date:     date,
	}
}
