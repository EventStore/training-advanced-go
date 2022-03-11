package doctorday

import (
	"testing"
	"time"

	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDayAggregate(t *testing.T) {
	doctorId := NewDoctorID(uuid.New())
	date := time.Date(2020, 5, 2, 10, 0, 0, 0, time.Local)
	store := infrastructure.NewFakeAggregateStore()
	registry := NewEventStoreDayRepository(store)
	a := DayTests{
		AggregateTests: infrastructure.NewAggregateTests(store),

		doctorId:   doctorId,
		patientId:  NewPatientID("John Doe"),
		date:       date,
		dayId:      NewDayID(doctorId, date),
		tenMinutes: time.Minute * time.Duration(10),
	}

	a.RegisterHandlers(NewHandlers(registry))

	t.Run("ShouldBeScheduled", a.ShouldBeScheduled)
	t.Run("ShouldNotBeScheduledTwice", a.ShouldNotBeScheduledTwice)
	t.Run("ShouldAllowToBookSlot", a.ShouldAllowToBookSlot)
	t.Run("ShouldNotAllowToBookSlotTwice", a.ShouldNotAllowToBookSlotTwice)
	t.Run("ShouldNotAllowToBookSlotIfDayNotScheduled", a.ShouldNotAllowToBookSlotIfDayNotScheduled)
	t.Run("ShouldNotAllowToBookAnUnscheduledSlot", a.ShouldNotAllowToBookAnUnscheduledSlot)
	t.Run("AllowToCancelBooking", a.AllowToCancelBooking)
	t.Run("NotAllowToCancelUnbookedSlot", a.NotAllowToCancelUnbookedSlot)
	t.Run("AllowToScheduleAnExtraSlot", a.AllowToScheduleAnExtraSlot)
	t.Run("DontAllowSchedulingOverlappingSlots", a.DontAllowSchedulingOverlappingSlots)
	t.Run("AllowToScheduleAdjacentSlots", a.AllowToScheduleAdjacentSlots)
	t.Run("CancelBookedSlotsWhenDayIsCancelled", a.CancelBookedSlotsWhenDayIsCancelled)
	t.Run("ArchiveScheduledDay", a.ArchiveScheduledDay)
}

func (t *DayTests) ShouldBeScheduled(tt *testing.T) {
	var slots []commands.ScheduledSlot
	slots = make([]commands.ScheduledSlot, 30)
	tenMinutes := time.Minute * 10
	for i, _ := range slots {
		slots[i] = commands.NewScheduledSlot(t.date.Add(tenMinutes*time.Duration(i)), tenMinutes)
	}

	t.Given()
	t.When(commands.NewScheduleDay(t.doctorId.Value, t.date, slots))
	t.Then(func(changes []interface{}, err error) {
		assert.NoError(tt, err)
		assert.Equal(tt, events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date), changes[0])
		assert.Len(tt, changes, 31)
	})
}

func (t *DayTests) ShouldNotBeScheduledTwice(tt *testing.T) {
	var slots []commands.ScheduledSlot
	slots = make([]commands.ScheduledSlot, 30)
	tenMinutes := time.Minute * 10
	for i, _ := range slots {
		slots[i] = commands.NewScheduledSlot(t.date.Add(tenMinutes*time.Duration(i)), tenMinutes)
	}

	t.Given(events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date))
	t.When(commands.NewScheduleDay(t.doctorId.Value, t.date, slots))
	t.ThenExpectError(tt, &DayAlreadyScheduledError{})
}

func (t *DayTests) ShouldAllowToBookSlot(tt *testing.T) {
	slotId := NewSlotID(uuid.New())
	expected := events.NewSlotBooked(slotId.Value, t.dayId.Value, "John Doe")

	t.Given(
		events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date),
		events.NewSlotScheduled(slotId.Value, t.dayId.Value, t.date, time.Minute*10))
	t.When(commands.NewBookSlot(slotId.Value, t.dayId.Value, "John Doe"))
	t.ThenExpectSingleChange(tt, expected)
}

func (t *DayTests) ShouldNotAllowToBookSlotTwice(tt *testing.T) {
	slotId := NewSlotID(uuid.New())

	t.Given(
		events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date),
		events.NewSlotScheduled(slotId.Value, t.dayId.Value, t.date, time.Minute*10),
		events.NewSlotBooked(slotId.Value, t.dayId.Value, "John Doe"))
	t.When(commands.NewBookSlot(slotId.Value, t.dayId.Value, "John Doe"))
	t.ThenExpectError(tt, &SlotAlreadyBookedError{})
}

func (t *DayTests) ShouldNotAllowToBookSlotIfDayNotScheduled(tt *testing.T) {
	slotId := NewSlotID(uuid.New())

	t.Given()
	t.When(commands.NewBookSlot(slotId.Value, t.dayId.Value, "John Doe"))
	t.ThenExpectError(tt, &DayNotScheduledError{})
}

func (t *DayTests) ShouldNotAllowToBookAnUnscheduledSlot(tt *testing.T) {
	slotId := NewSlotID(uuid.New())

	t.Given(events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date))
	t.When(commands.NewBookSlot(slotId.Value, t.dayId.Value, "John Doe"))
	t.ThenExpectError(tt, &SlotNotScheduledError{})
}

func (t *DayTests) AllowToCancelBooking(tt *testing.T) {
	slotId := NewSlotID(uuid.New())
	reason := "Cancel reason"
	expected := events.NewSlotBookingCancelled(slotId.Value, t.dayId.Value, reason)

	t.Given(
		events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date),
		events.NewSlotScheduled(slotId.Value, t.dayId.Value, t.date, time.Minute*10),
		events.NewSlotBooked(slotId.Value, t.dayId.Value, "John Doe"))
	t.When(commands.NewCancelSlotBooking(slotId.Value, t.dayId.Value, reason))
	t.ThenExpectSingleChange(tt, expected)
}

func (t *DayTests) NotAllowToCancelUnbookedSlot(tt *testing.T) {
	slotId := NewSlotID(uuid.New())

	t.Given(
		events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date),
		events.NewSlotScheduled(slotId.Value, t.dayId.Value, t.date, time.Minute*10))
	t.When(commands.NewCancelSlotBooking(slotId.Value, t.dayId.Value, "Some reason"))
	t.ThenExpectError(tt, &SlotNotBookedError{})
}

func (t *DayTests) AllowToScheduleAnExtraSlot(tt *testing.T) {
	slotId := NewSlotID(uuid.New())
	expected := events.NewSlotScheduled(slotId.Value, t.dayId.Value, t.date, t.tenMinutes)

	t.Given(
		events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date),
		/* events.NewSlotScheduled(slotId.Value, t.dayId.Value, t.date, time.Minute*10) */)
	t.When(commands.NewScheduleSlot(slotId.Value, t.doctorId.Value, t.date, t.tenMinutes))
	t.ThenExpectSingleChange(tt, expected)
}

func (t *DayTests) DontAllowSchedulingOverlappingSlots(tt *testing.T) {
	slotOneId := NewSlotID(uuid.New())
	slotTwoId := NewSlotID(uuid.New())

	t.Given(
		events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date),
		events.NewSlotScheduled(slotOneId.Value, t.dayId.Value, t.date, t.tenMinutes))
	t.When(commands.NewScheduleSlot(slotTwoId.Value, t.doctorId.Value, t.date, t.tenMinutes))
	t.ThenExpectError(tt, &SlotOverlappedError{})
}

func (t *DayTests) AllowToScheduleAdjacentSlots(tt *testing.T) {
	slotOneId := NewSlotID(uuid.New())
	slotTwoId := NewSlotID(uuid.New())
	expected := events.NewSlotScheduled(slotTwoId.Value, t.dayId.Value, t.date.Add(t.tenMinutes), t.tenMinutes)

	t.Given(
		events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date),
		events.NewSlotScheduled(slotOneId.Value, t.dayId.Value, t.date, t.tenMinutes))
	t.When(commands.NewScheduleSlot(slotTwoId.Value, t.doctorId.Value, t.date.Add(t.tenMinutes), t.tenMinutes))
	t.ThenExpectSingleChange(tt, expected)
}

func (t *DayTests) CancelBookedSlotsWhenDayIsCancelled(tt *testing.T) {
	slotOneId := NewSlotID(uuid.New())
	slotTwoId := NewSlotID(uuid.New())

	t.Given(
		events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date),
		events.NewSlotScheduled(slotOneId.Value, t.dayId.Value, t.date, t.tenMinutes),
		events.NewSlotScheduled(slotTwoId.Value, t.dayId.Value, t.date.Add(t.tenMinutes), t.tenMinutes),
		events.NewSlotBooked(slotOneId.Value, t.dayId.Value, t.patientId.Value))
	t.When(commands.NewCancelDaySchedule(t.dayId.Value))
	t.ThenExpectChanges(tt, []interface{}{
		events.NewSlotBookingCancelled(slotOneId.Value, t.dayId.Value, DayCancelledReason),
		events.NewSlotScheduleCancelled(slotOneId.Value, t.dayId.Value),
		events.NewSlotScheduleCancelled(slotTwoId.Value, t.dayId.Value),
		events.NewDayScheduleCancelled(t.dayId.Value),
	})
}

func (t *DayTests) ArchiveScheduledDay(tt *testing.T) {
	expected := events.NewDayScheduleArchived(t.dayId.Value)

	t.Given(events.NewDayScheduled(t.dayId.Value, t.doctorId.Value, t.date))
	t.When(commands.NewArchiveDaySchedule(t.dayId.Value))
	t.ThenExpectSingleChange(tt, expected)
}

type DayTests struct {
	infrastructure.AggregateTests

	dayId      DayID
	doctorId   DoctorID
	patientId  PatientID
	date       time.Time
	tenMinutes time.Duration
}
