package application

import (
	"context"
	"testing"
	"time"

	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/EventStore/training-introduction-go/infrastructure/mongodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAvailableSlotsHandler(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost"))
	defer client.Disconnect(context.TODO())
	assert.NoError(t, err)

	p := &AvailableSlotsTests{
		HandlerTests: infrastructure.NewHandlerTests(t),

		dayId:      "dayId",
		patientId:  "patientId",
		reason:     "Some cancellation reason",
		slotId:     uuid.New(),
		now:        time.Now(),
		tenMinutes: 10 * time.Minute,
	}

	p.SetHandlerFactory(func() infrastructure.EventHandler {
		p.repository = mongodb.NewAvailableSlotsRepository(client.Database(uuid.NewString()))
		return NewAvailableSlotsProjection(p.repository)
	})

	t.Run("ShouldAddSlotToTheList", p.ShouldAddSlotToTheList)
	t.Run("ShouldHideTheSlotFromListIfBooked", p.ShouldHideTheSlotFromListIfBooked)
	t.Run("ShouldShowSlotIfBookingWasCancelled", p.ShouldShowSlotIfBookingWasCancelled)
	t.Run("ShouldDeleteSlotIfSlotWasCancelled", p.ShouldDeleteSlotIfSlotWasCancelled)
}

func (p *AvailableSlotsTests) ShouldAddSlotToTheList(t *testing.T) {
	p.Given(events.NewSlotScheduled(p.slotId, p.dayId, p.now, p.tenMinutes))
	p.Then(
		[]readmodel.AvailableSlot{
			readmodel.NewAvailableSlot(
				p.slotId.String(),
				p.dayId,
				p.now.Format("02-01-2006"),
				p.now.Format("15:04:05"),
				p.tenMinutes)},
		p.getSlotsAvailableOn(p.now))
}

func (p *AvailableSlotsTests) ShouldHideTheSlotFromListIfBooked(t *testing.T) {
	p.Given(
		events.NewSlotScheduled(p.slotId, p.dayId, p.now, p.tenMinutes),
		events.NewSlotBooked(p.slotId, p.dayId, p.patientId))
	p.Then(
		[]readmodel.AvailableSlot{},
		p.getSlotsAvailableOn(p.now))
}

func (p *AvailableSlotsTests) ShouldShowSlotIfBookingWasCancelled(t *testing.T) {
	p.Given(
		events.NewSlotScheduled(p.slotId, p.dayId, p.now, p.tenMinutes),
		events.NewSlotBooked(p.slotId, p.dayId, p.patientId),
		events.NewSlotBookingCancelled(p.slotId, p.dayId, p.reason))
	p.Then(
		[]readmodel.AvailableSlot{
			readmodel.NewAvailableSlot(
				p.slotId.String(),
				p.dayId,
				p.now.Format("02-01-2006"),
				p.now.Format("15:04:05"),
				p.tenMinutes)},
		p.getSlotsAvailableOn(p.now))
}

func (p *AvailableSlotsTests) ShouldDeleteSlotIfSlotWasCancelled(t *testing.T) {
	p.Given(
		events.NewSlotScheduled(p.slotId, p.dayId, p.now, p.tenMinutes),
		events.NewSlotScheduleCancelled(p.slotId, p.dayId))
	p.Then(
		[]readmodel.AvailableSlot{},
		p.getSlotsAvailableOn(p.now))
}

func (p *AvailableSlotsTests) getSlotsAvailableOn(now time.Time) []readmodel.AvailableSlot {
	result, err := p.repository.GetSlotsAvailableOn(p.now)
	assert.NoError(p.T, err)

	return result
}

type AvailableSlotsTests struct {
	infrastructure.HandlerTests

	repository *mongodb.AvailableSlotsRepository

	dayId      string
	patientId  string
	reason     string
	slotId     uuid.UUID
	now        time.Time
	tenMinutes time.Duration
}
