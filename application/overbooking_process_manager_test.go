package application

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/EventStore/training-introduction-go/domain/doctorday"
	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/eventsourcing"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/EventStore/training-introduction-go/infrastructure/mongodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestOverbooking(t *testing.T) {
	client, err := createESDBClient()
	assert.NoError(t, err)

	tm := eventsourcing.NewTypeMapper()
	doctorday.RegisterTypes(tm)
	esSerde := infrastructure.NewEsEventSerde(tm)
	tenantPrefix := fmt.Sprintf("overbooking_tests_%s_", uuid.NewString())
	eventStore := infrastructure.NewEsEventStore(client, tenantPrefix, esSerde)
	cmdStore := infrastructure.NewEsCommandStore(eventStore, nil, nil, nil)

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost"))
	defer mongoClient.Disconnect(context.TODO())
	assert.NoError(t, err)

	p := &OverbookingTests{
		HandlerTests: infrastructure.NewHandlerTests(t),

		bookingLimitPerPatient: 3,

		dayId:      "dayId",
		patientId:  "patient 1",
		reason:     "Some cancellation reason",
		slotId:     uuid.New(),
		now:        time.Now().Truncate(time.Second),
		tenMinutes: 10 * time.Minute,
		eventStore: eventStore,
	}

	p.SetHandlerFactory(func() infrastructure.EventHandler {
		p.repository = mongodb.NewBookedSlotsRepository(mongoClient.Database(uuid.NewString()))
		return NewOverbookingProcessManager(p.repository, cmdStore, p.bookingLimitPerPatient)
	})

	t.Run("ShouldIncrementTheVisitCounterWhenSlotIsBooked", p.ShouldIncrementTheVisitCounterWhenSlotIsBooked)
	t.Run("ShouldDecrementTheVisitCounterWhenSlotBookingIsCancelled", p.ShouldDecrementTheVisitCounterWhenSlotBookingIsCancelled)
	t.Run("ShouldIssueCommandToCancelSlotIfBookingLimitWasReached", p.ShouldIssueCommandToCancelSlotIfBookingLimitWasReached)
}

func (a *OverbookingTests) ShouldIncrementTheVisitCounterWhenSlotIsBooked(t *testing.T) {
	dayId := uuid.NewString()
	slotSchedule1 := events.NewSlotScheduled(uuid.New(), dayId, a.now, a.tenMinutes)
	slotSchedule2 := events.NewSlotScheduled(uuid.New(), dayId, a.now.Add(10*time.Minute), a.tenMinutes)
	slotBooked1 := events.NewSlotBooked(slotSchedule1.SlotId, dayId, a.patientId)
	slotBooked2 := events.NewSlotBooked(slotSchedule2.SlotId, dayId, a.patientId)

	a.Given(slotSchedule1, slotSchedule2, slotBooked1, slotBooked2)

	count, err := a.repository.CountByPatientAndMonth(a.patientId, int(a.now.Month()))
	assert.NoError(t, err)

	a.Then(2, count)
}

func (a *OverbookingTests) ShouldDecrementTheVisitCounterWhenSlotBookingIsCancelled(t *testing.T) {
	dayId := uuid.NewString()
	slotSchedule1 := events.NewSlotScheduled(uuid.New(), dayId, a.now, a.tenMinutes)
	slotSchedule2 := events.NewSlotScheduled(uuid.New(), dayId, a.now.Add(10*time.Minute), a.tenMinutes)
	slotBooked1 := events.NewSlotBooked(slotSchedule1.SlotId, dayId, a.patientId)
	slotBooked2 := events.NewSlotBooked(slotSchedule2.SlotId, dayId, a.patientId)
	slotBookingCancelled := events.NewSlotBookingCancelled(slotSchedule2.SlotId, dayId, "no longer needed")

	a.Given(slotSchedule1, slotSchedule2, slotBooked1, slotBooked2, slotBookingCancelled)

	count, err := a.repository.CountByPatientAndMonth(a.patientId, int(a.now.Month()))
	assert.NoError(t, err)

	a.Then(1, count)
}

func (a *OverbookingTests) ShouldIssueCommandToCancelSlotIfBookingLimitWasReached(t *testing.T) {
	dayId := uuid.NewString()
	slotSchedule1 := events.NewSlotScheduled(uuid.New(), dayId, a.now, a.tenMinutes)
	slotSchedule2 := events.NewSlotScheduled(uuid.New(), dayId, a.now.Add(10*time.Minute), a.tenMinutes)
	slotSchedule3 := events.NewSlotScheduled(uuid.New(), dayId, a.now.Add(20*time.Minute), a.tenMinutes)
	slotSchedule4 := events.NewSlotScheduled(uuid.New(), dayId, a.now.Add(30*time.Minute), a.tenMinutes)
	slotBooked1 := events.NewSlotBooked(slotSchedule1.SlotId, dayId, a.patientId)
	slotBooked2 := events.NewSlotBooked(slotSchedule2.SlotId, dayId, a.patientId)
	slotBooked3 := events.NewSlotBooked(slotSchedule3.SlotId, dayId, a.patientId)
	slotBooked4 := events.NewSlotBooked(slotSchedule4.SlotId, dayId, a.patientId)

	a.Given(
		slotSchedule1, slotSchedule2, slotSchedule3, slotSchedule4,
		slotBooked1, slotBooked2, slotBooked3, slotBooked4)

	cmd, err := a.eventStore.LoadCommand("async_command_handler-day")
	assert.NoError(t, err)

	a.Then(
		commands.NewCancelSlotBooking(slotSchedule4.SlotId, dayId, "overbooked"),
		cmd[0].Command)
}

type OverbookingTests struct {
	infrastructure.HandlerTests

	eventStore *infrastructure.EsEventStore
	repository readmodel.BookedSlotsRepository

	bookingLimitPerPatient int

	dayId      string
	patientId  string
	reason     string
	slotId     uuid.UUID
	now        time.Time
	tenMinutes time.Duration
}
