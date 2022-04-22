package application

import (
	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/google/uuid"
)

type OverbookingProcessManager struct {
	infrastructure.EventHandlerBase
}

func NewOverbookingProcessManager(r readmodel.BookedSlotsRepository, c infrastructure.CommandStore, bookingLimitPerPatient int) *OverbookingProcessManager {
	h := infrastructure.NewEventHandler()

	h.When(events.SlotScheduled{}, func(e interface{}, _ infrastructure.EventMetadata) error {
		s := e.(events.SlotScheduled)
		r.AddSlot(readmodel.NewBookedSlot(s.SlotId.String(), s.DayId, int(s.StartTime.Month())))
		return nil
	})

	h.When(events.SlotBooked{}, func(e interface{}, m infrastructure.EventMetadata) error {
		s := e.(events.SlotBooked)
		r.MarkSlotAsBooked(s.SlotId.String(), s.PatientId)

		slot, err := r.GetSlot(s.SlotId.String())
		if err != nil {
			return err
		}
		count, err := r.CountByPatientAndMonth(s.PatientId, slot.Month)
		if err != nil {
			return err
		}
		if count > bookingLimitPerPatient {
			metadata := infrastructure.NewCommandMetadata(m.CorrelationId.Value, uuid.New())
			err := c.Send(commands.NewCancelSlotBooking(s.SlotId, slot.DayId, "overbooked"), metadata)
			if err != nil {
				return err
			}
		}
		return nil
	})

	h.When(events.SlotBookingCancelled{}, func(e interface{}, _ infrastructure.EventMetadata) error {
		s := e.(events.SlotBookingCancelled)
		r.MarkSlotAsAvailable(s.SlotId.String())
		return nil
	})

	return &OverbookingProcessManager{h}
}
