package application

import (
	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/EventStore/training-introduction-go/infrastructure/mongodb"
)

type AvailableSlotsProjection struct {
	infrastructure.EventHandlerBase
}

func NewAvailableSlotsProjection(r *mongodb.AvailableSlotsRepository) *AvailableSlotsProjection {
	h := infrastructure.NewEventHandler()
	h.When(events.SlotScheduled{}, func(e interface{}, _ infrastructure.EventMetadata) error {
		s := e.(events.SlotScheduled)
		return r.AddSlot(mongodb.AvailableSlot{
			Id:        s.SlotId.String(),
			DayId:     s.DayId,
			Date:      s.StartTime.Format("02-01-2006"),
			StartTime: s.StartTime.Format("15:04:05"),
			Duration:  s.Duration})
	})

	h.When(events.SlotBooked{}, func(e interface{}, _ infrastructure.EventMetadata) error {
		b := e.(events.SlotBooked)
		return r.HideSlot(b.SlotId)
	})

	h.When(events.SlotBookingCancelled{}, func(e interface{}, _ infrastructure.EventMetadata) error {
		c := e.(events.SlotBookingCancelled)
		return r.ShowSlot(c.SlotId)
	})

	h.When(events.SlotScheduleCancelled{}, func(e interface{}, _ infrastructure.EventMetadata) error {
		c := e.(events.SlotScheduleCancelled)
		return r.DeleteSlot(c.SlotId)
	})

	return &AvailableSlotsProjection{h}
}
