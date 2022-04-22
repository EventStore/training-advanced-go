package readmodel

import "time"

type AvailableSlotsRepository interface {
	GetSlotsAvailableOn(time time.Time) ([]AvailableSlot, error)
}
