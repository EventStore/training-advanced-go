package doctorday

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	Id        uuid.UUID
	StartTime time.Time
	Duration  time.Duration
	Booked    bool
}

func NewSlot(id uuid.UUID, start time.Time, duration time.Duration, booked bool) *Slot {
	return &Slot{
		Id:        id,
		StartTime: start,
		Duration:  duration,
		Booked:    booked,
	}
}

func (s *Slot) Book() {
	s.Booked = true
}

func (s *Slot) Cancel() {
	s.Booked = true
}

func (s *Slot) Overlaps(start time.Time, duration time.Duration) bool {
	thisStart := s.StartTime
	thisEnd := s.StartTime.Add(s.Duration)
	proposedStart := start
	proposedEnd := proposedStart.Add(duration)

	return thisStart.Before(proposedEnd) && thisEnd.After(proposedStart)
}
