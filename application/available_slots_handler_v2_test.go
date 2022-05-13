package application

import (
	"context"
	"encoding/json"
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

func TestAvailableSlotsHandlerV2(t *testing.T) {
	s := events.NewSlotScheduled(uuid.New(), "dayId", time.Now(), 10*time.Minute)
	m := map[string]interface{}{}

	b, e := json.Marshal(s)

	e = json.Unmarshal(b, &m)
	assert.NoError(t, e)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost"))
	defer client.Disconnect(context.TODO())
	assert.NoError(t, err)

	p := &AvailableSlotsV2Tests{
		HandlerTests: infrastructure.NewHandlerTests(t),

		dayId:      "dayId",
		patientId:  "patientId",
		reason:     "Some cancellation reason",
		slotId:     uuid.New(),
		now:        time.Now(),
		tenMinutes: 10 * time.Minute,
	}

	// Repeats every event 2x, e.g.: 1 1 2 2 3 3
	p.EnableAtLeastOnceMonkey = false
	// Repeats all elements except last e.g.: 1 2 3 1 2
	//p.EnableAtLeastOnceGorilla = false

	p.SetHandlerFactory(func() infrastructure.EventHandler {
		p.repository = mongodb.NewAvailableSlotsRepository(client.Database(uuid.NewString()))
		return NewAvailableSlotsProjectionV2(p.repository)
	})

	t.Run("ShouldAddSlotToTheList", p.ShouldAddSlotToTheList)
	t.Run("ShouldHideTheSlotFromListIfBooked", p.ShouldHideTheSlotFromListIfBooked)
	t.Run("ShouldShowSlotIfBookingWasCancelled", p.ShouldShowSlotIfBookingWasCancelled)
	t.Run("ShouldDeleteSlotIfSlotWasCancelled", p.ShouldDeleteSlotIfSlotWasCancelled)
}

func (p *AvailableSlotsV2Tests) ShouldAddSlotToTheList(t *testing.T) {
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

func (p *AvailableSlotsV2Tests) ShouldHideTheSlotFromListIfBooked(t *testing.T) {
	p.Given(
		events.NewSlotScheduled(p.slotId, p.dayId, p.now, p.tenMinutes),
		events.NewSlotBooked(p.slotId, p.dayId, p.patientId))
	p.Then(
		[]readmodel.AvailableSlot{},
		p.getSlotsAvailableOn(p.now))
}

func (p *AvailableSlotsV2Tests) ShouldShowSlotIfBookingWasCancelled(t *testing.T) {
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

func (p *AvailableSlotsV2Tests) ShouldDeleteSlotIfSlotWasCancelled(t *testing.T) {
	p.Given(
		events.NewSlotScheduled(p.slotId, p.dayId, p.now, p.tenMinutes),
		events.NewSlotScheduleCancelled(p.slotId, p.dayId))
	p.Then(
		[]readmodel.AvailableSlot{},
		p.getSlotsAvailableOn(p.now))
}

func (p *AvailableSlotsV2Tests) getSlotsAvailableOn(now time.Time) []readmodel.AvailableSlot {
	result, err := p.repository.GetSlotsAvailableOn(p.now)
	assert.NoError(p.T, err)

	return result
}

type AvailableSlotsV2Tests struct {
	infrastructure.HandlerTests

	repository *mongodb.AvailableSlotsRepository

	dayId      string
	patientId  string
	reason     string
	slotId     uuid.UUID
	now        time.Time
	tenMinutes time.Duration
}
