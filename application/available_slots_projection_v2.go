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
		return nil
	})

	p.When(events.SlotBooked{}, func(e interface{}, m infrastructure.EventMetadata) error {
		return nil
	})

	// Add when for SlotBookingCancelled

	return &AvailableSlotsProjectionV2{p}
}
