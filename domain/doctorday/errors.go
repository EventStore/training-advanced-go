package doctorday

import "fmt"

type DayAlreadyScheduledError struct{}

func (e DayAlreadyScheduledError) Error() string {
	return fmt.Sprintf("day already scheduled error")
}

type DayNotScheduledError struct{}

func (e DayNotScheduledError) Error() string {
	return fmt.Sprintf("day not scheduled error")
}

type DayScheduleAlreadyArchivedError struct{}

func (e DayScheduleAlreadyArchivedError) Error() string {
	return fmt.Sprintf("day schedule already archived error")
}

type DayScheduleAlreadyCancelledError struct{}

func (e DayScheduleAlreadyCancelledError) Error() string {
	return fmt.Sprintf("day schedule already cancelled error")
}

type SlotAlreadyBookedError struct{}

func (e SlotAlreadyBookedError) Error() string {
	return fmt.Sprintf("slot already booked error")
}

type SlotNotBookedError struct{}

func (e SlotNotBookedError) Error() string {
	return fmt.Sprintf("slot not booked error")
}

type SlotNotScheduledError struct{}

func (e SlotNotScheduledError) Error() string {
	return fmt.Sprintf("slot not scheduled error")
}

type SlotOverlappedError struct{}

func (e SlotOverlappedError) Error() string {
	return fmt.Sprintf("slot overlapped error")
}
