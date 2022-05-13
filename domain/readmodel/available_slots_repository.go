package readmodel

import "time"

type AvailableSlotsRepository interface {
	GetSlotsAvailableOn(time time.Time) ([]AvailableSlot, error)

	AddSlot(slot AvailableSlot) error

	HideSlot(slotId string) error

	ShowSlot(slotId string) error

	DeleteSlot(slotId string) error
}
