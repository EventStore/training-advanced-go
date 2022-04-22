package application

import (
	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/EventStore/training-introduction-go/infrastructure/mongodb"
)

type AvailableSlotsProjectionV2 struct {
	infrastructure.EventHandlerBase
}

func NewAvailableSlotsProjectionV2(r *mongodb.AvailableSlotsRepositoryV2) *AvailableSlotsProjectionV2 {
	p := infrastructure.NewEventHandler()
	p.When(events.SlotScheduled{}, func(e interface{}, m infrastructure.EventMetadata) error {
		s := e.(events.SlotScheduled)
		return r.AddSlot(mongodb.AvailableSlotV2{
			Id:        s.SlotId.String(),
			DayId:     s.DayId,
			Date:      s.StartTime.Format("02-01-2006"),
			StartTime: s.StartTime.Format("15:04:05"),
			Duration:  s.Duration})
	})

	p.When(events.SlotBooked{}, func(e interface{}, m infrastructure.EventMetadata) error {
		b := e.(events.SlotBooked)
		return r.HideSlot(b.SlotId)
	})

	p.When(events.SlotBookingCancelled{}, func(e interface{}, m infrastructure.EventMetadata) error {
		c := e.(events.SlotBookingCancelled)
		return r.ShowSlot(c.SlotId)
	})

	p.When(events.SlotScheduleCancelled{}, func(e interface{}, m infrastructure.EventMetadata) error {
		c := e.(events.SlotScheduleCancelled)
		return r.DeleteSlot(c.SlotId)
	})

	return &AvailableSlotsProjectionV2{p}
}
