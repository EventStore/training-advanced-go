package application

import (
	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/infrastructure"
)

type OverbookingProcessManager struct {
	infrastructure.EventHandlerBase
}

func NewOverbookingProcessManager(r readmodel.BookedSlotsRepository, c infrastructure.CommandStore, bookingLimitPerPatient int) *OverbookingProcessManager {
	h := infrastructure.NewEventHandler()

	h.When(events.SlotScheduled{}, func(e interface{}, _ infrastructure.EventMetadata) error {
		return nil
	})

	h.When(events.SlotBooked{}, func(e interface{}, m infrastructure.EventMetadata) error {
		return nil
	})

	h.When(events.SlotBookingCancelled{}, func(e interface{}, _ infrastructure.EventMetadata) error {
		return nil
	})

	return &OverbookingProcessManager{h}
}
