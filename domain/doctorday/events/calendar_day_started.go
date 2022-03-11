package events

import "time"

type CalendarDayStarted struct {
	Date time.Time
}

func NewCalendarDayStarted(date time.Time) CalendarDayStarted {
	return CalendarDayStarted{
		Date: date,
	}
}
