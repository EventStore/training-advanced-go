package doctorday
//
//import (
//	"testing"
//
//	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
//	"github.com/EventStore/training-introduction-go/eventsourcing"
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestCheckSlotBookingCancelledCorrectlyMapsWithDefaultValue(t *testing.T) {
//	typeMapper := eventsourcing.NewTypeMapper()
//	RegisterTypes(typeMapper)
//
//	dataToType, err := typeMapper.GetDataToType("doctorday-slot-booking-cancelled")
//	assert.NoError(t, err)
//
//	slotId := uuid.New()
//	slotBookingCancelled := dataToType(map[string]interface{}{
//		"dayId":  "dayId",
//		"slotId": slotId.String(),
//		"reason": "reason",
//	})
//
//	assert.NotNil(t, slotBookingCancelled)
//	assert.IsType(t, events.SlotBookingCancelled{}, slotBookingCancelled)
//	assert.Equal(t, events.NewSlotBookingCancelled(slotId, "dayId", "reason", "unknown request"), slotBookingCancelled)
//}
//
//func TestCheckSlotBookingCancelledCorrectlyMapsWithValuePresent(t *testing.T) {
//	typeMapper := eventsourcing.NewTypeMapper()
//	RegisterTypes(typeMapper)
//
//	dataToType, err := typeMapper.GetDataToType("doctorday-slot-booking-cancelled")
//	assert.NoError(t, err)
//
//	slotId := uuid.New()
//	slotBookingCancelled := dataToType(map[string]interface{}{
//		"dayId":       "dayId",
//		"slotId":      slotId.String(),
//		"reason":      "reason",
//		"requestedBy": "doctor",
//	})
//
//	assert.NotNil(t, slotBookingCancelled)
//	assert.IsType(t, events.SlotBookingCancelled{}, slotBookingCancelled)
//	assert.Equal(t, events.NewSlotBookingCancelled(slotId, "dayId", "reason", "doctor"), slotBookingCancelled)
//}
