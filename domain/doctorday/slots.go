package doctorday

import (
	"time"

	"github.com/google/uuid"
)

type Slots struct {
	slots []*Slot
}

func (s *Slots) Add(id uuid.UUID, start time.Time, duration time.Duration, booked bool) {
	s.slots = append(s.slots, NewSlot(id, start, duration, booked))
}

func (s *Slots) Remove(id SlotID) {
	slots := make([]*Slot, 0)
	for _, slot := range s.slots {
		if slot.Id != id.Value {
			slots = append(slots, slot)
		}
	}

	s.slots = slots
}

func (s *Slots) Overlaps(start time.Time, duration time.Duration) bool {
	for _, slot := range s.slots {
		if slot.Overlaps(start, duration) {
			return true
		}
	}

	return false
}

func (s *Slots) GetStatus(id SlotID) SlotStatus {
	slot := s.getSlot(id)
	if slot == nil {
		return SlotNotScheduled
	}

	if slot.Booked {
		return SlotBooked
	}

	return SlotAvailable
}

func (s *Slots) MarkAsBooked(id SlotID) {
	slot := s.getSlot(id)
	if slot != nil {
		slot.Book()
	}
}

func (s *Slots) MarkAsAvailable(id SlotID) {
	slot := s.getSlot(id)
	if slot != nil {
		slot.Cancel()
	}
}

func (s *Slots) HasBookedSlot(id SlotID) bool {
	slot := s.getSlot(id)
	if slot == nil {
		return false
	}

	return slot.Booked
}

func (s *Slots) GetBookedSlots() []*Slot {
	bookedSlots := make([]*Slot, 0)
	for _, slot := range s.slots {
		if slot.Booked {
			bookedSlots = append(bookedSlots, slot)
		}
	}

	return bookedSlots
}

func (s *Slots) GetAllSlots() []*Slot {
	return s.slots
}

func (s *Slots) getSlot(id SlotID) *Slot {
	for _, slot := range s.slots {
		if slot.Id == id.Value {
			return slot
		}
	}

	return nil
}
