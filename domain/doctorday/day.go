package doctorday

import (
	"fmt"
	"time"

	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/eventsourcing"
	"github.com/google/uuid"
)

const (
	DayCancelledReason = "day cancelled"
)

type Day struct {
	eventsourcing.AggregateRootSnapshotBase

	isArchived  bool
	isCancelled bool
	isScheduled bool
	slots       Slots
}

func NewDay() *Day {
	a := &Day{
		AggregateRootSnapshotBase: eventsourcing.NewAggregateRootSnapshot(),
	}

	a.Register(events.DayScheduled{}, func(e interface{}) { a.DayScheduled(e.(events.DayScheduled)) })
	a.Register(events.SlotScheduled{}, func(e interface{}) { a.SlotScheduled(e.(events.SlotScheduled)) })
	a.Register(events.SlotBooked{}, func(e interface{}) { a.SlotBooked(e.(events.SlotBooked)) })
	a.Register(events.SlotBookingCancelled{}, func(e interface{}) { a.SlotBookingCancelled(e.(events.SlotBookingCancelled)) })
	a.Register(events.SlotScheduleCancelled{}, func(e interface{}) { a.SlotScheduleCancelled(e.(events.SlotScheduleCancelled)) })
	a.Register(events.DayScheduleArchived{}, func(e interface{}) { a.DayScheduleArchived(e.(events.DayScheduleArchived)) })
	a.RegisterSnapshot(func(s interface{}) { a.loadSnapshot(s.(DaySnapshot)) }, a.getSnapshot)

	return a
}

// Schedule day

func (s *Day) ScheduleDay(doctorId DoctorID, date time.Time, slots []commands.ScheduledSlot) error {
	err := s.isDayCancelledOrArchived()
	if err != nil {
		return err
	}

	if s.isScheduled {
		return &DayAlreadyScheduledError{}
	}

	dayId := NewDayID(doctorId, date)
	s.Raise(events.NewDayScheduled(dayId.Value, doctorId.Value, date))

	for _, slot := range slots {
		s.Raise(events.NewSlotScheduled(uuid.New(), dayId.Value, slot.StartTime, slot.Duration))
	}

	return nil
}

func (s *Day) DayScheduled(e events.DayScheduled) {
	s.Id = NewDayID(NewDoctorID(e.DoctorId), e.Date).Value
	s.isScheduled = true
}

// Schedule slot

func (s *Day) ScheduleSlot(slotId uuid.UUID, start time.Time, duration time.Duration) error {
	err := s.isDayCancelledOrArchived()
	if err != nil {
		return err
	}
	err = s.isDayNotScheduled()
	if err != nil {
		return err
	}

	if s.slots.Overlaps(start, duration) {
		return &SlotOverlappedError{}
	}

	s.Raise(events.NewSlotScheduled(slotId, s.Id, start, duration))
	return nil
}

func (s *Day) SlotScheduled(e events.SlotScheduled) {
	s.slots.Add(e.SlotId, e.StartTime, e.Duration, false)
}

// Book slot

func (s *Day) BookSlot(slotId SlotID, patientId PatientID) error {
	err := s.isDayCancelledOrArchived()
	if err != nil {
		return err
	}
	err = s.isDayNotScheduled()
	if err != nil {
		return err
	}

	slotStatus := s.slots.GetStatus(slotId)

	switch slotStatus {
	case SlotAvailable:
		s.Raise(events.NewSlotBooked(slotId.Value, s.Id, patientId.Value))
		return nil
	case SlotBooked:
		return &SlotAlreadyBookedError{}
	case SlotNotScheduled:
		return &SlotNotScheduledError{}
	default:
		return fmt.Errorf("invalid slot status: %d", slotStatus)
	}

	return fmt.Errorf("unexpected slot booking error")
}

func (s *Day) SlotBooked(e events.SlotBooked) {
	s.slots.MarkAsBooked(NewSlotID(e.SlotId))
}

// Cancel slot booking

func (s *Day) CancelSlotBooking(slotId SlotID, reason string) error {
	err := s.isDayCancelledOrArchived()
	if err != nil {
		return err
	}
	err = s.isDayNotScheduled()
	if err != nil {
		return err
	}

	if !s.slots.HasBookedSlot(slotId) {
		return &SlotNotBookedError{}
	}

	s.Raise(events.NewSlotBookingCancelled(slotId.Value, s.Id, reason))
	return nil
}

func (s *Day) SlotBookingCancelled(e events.SlotBookingCancelled) {
	s.slots.MarkAsAvailable(NewSlotID(e.SlotId))
}

// Cancel day

func (s *Day) Cancel() error {
	return nil
}

func (s *Day) SlotScheduleCancelled(e events.SlotScheduleCancelled) {
	s.slots.Remove(NewSlotID(e.SlotId))
}

// Archive day

func (s *Day) Archive() error {
	err := s.isDayNotScheduled()
	if err != nil {
		return err
	}

	if s.isArchived {
		return &DayScheduleAlreadyArchivedError{}
	}

	s.Raise(events.NewDayScheduleArchived(s.Id))
	return nil
}

func (s *Day) DayScheduleArchived(_ events.DayScheduleArchived) {
	s.isArchived = true
}

func (s *Day) isDayCancelledOrArchived() error {
	if s.isArchived {
		return &DayScheduleAlreadyArchivedError{}
	}

	if s.isCancelled {
		return &DayScheduleAlreadyCancelledError{}
	}

	return nil
}

func (s *Day) isDayNotScheduled() error {
	if !s.isScheduled {
		return &DayNotScheduledError{}
	}

	return nil
}

// Snapshot

func (s *Day) getSnapshot() interface{} {
	slots := make([]SlotSnapshot, 0)
	for _, slot := range s.slots.GetAllSlots() {
		slots = append(slots, NewSlotSnapshot(slot.Id, slot.StartTime, slot.Duration, slot.Booked))
	}

	return NewDaySnapshot(s.isArchived, s.isCancelled, s.isScheduled, slots)
}

func (s *Day) loadSnapshot(snapshot DaySnapshot) {
	s.isArchived = snapshot.IsArchived
	s.isCancelled = snapshot.IsCancelled
	s.isScheduled = snapshot.IsScheduled

	for _, slot := range snapshot.Slots {
		s.slots.Add(slot.Id, slot.Start, slot.Duration, slot.Booked)
	}
}
