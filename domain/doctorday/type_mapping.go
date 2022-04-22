package doctorday

import (
	"time"

	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/eventsourcing"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/google/uuid"
)

const Prefix = "doctorday"

func RegisterTypes(tm *eventsourcing.TypeMapper) {
	mustParseDate := func(s string) time.Time {
		t, _ := time.Parse("2006-01-02", s)
		return t
	}
	mustParseTime := func(s string) time.Time {
		t, _ := time.Parse(time.RFC3339, s)
		return t
	}
	mustParseDuration := func(s string) time.Duration {
		d, _ := time.ParseDuration(s)
		return d
	}

	tm.MapEvent(infrastructure.GetValueType(events.DayScheduled{}), Prefix+"-day-scheduled",
		func(d map[string]interface{}) interface{} {
			return events.NewDayScheduled(
				d["dayId"].(string),
				uuid.MustParse(d["doctorId"].(string)),
				mustParseTime(d["date"].(string)))
		},
		func(v interface{}) map[string]interface{} {
			t := v.(events.DayScheduled)
			return map[string]interface{}{
				"dayId":    t.DayId,
				"doctorId": t.DoctorId,
				"date":     t.Date,
			}
		})

	tm.MapEvent(infrastructure.GetValueType(events.SlotScheduled{}), Prefix+"-slot-scheduled",
		func(d map[string]interface{}) interface{} {
			return events.NewSlotScheduled(
				uuid.MustParse(d["slotId"].(string)),
				d["dayId"].(string),
				mustParseTime(d["startTime"].(string)),
				mustParseDuration(d["duration"].(string)))
		},
		func(v interface{}) map[string]interface{} {
			t := v.(events.SlotScheduled)
			return map[string]interface{}{
				"slotId":    t.SlotId,
				"dayId":     t.DayId,
				"startTime": t.StartTime,
				"duration":  t.Duration.String(),
			}
		})

	tm.MapEvent(infrastructure.GetValueType(events.SlotBooked{}), Prefix+"-slot-booked",
		func(d map[string]interface{}) interface{} {
			return events.NewSlotBooked(
				uuid.MustParse(d["slotId"].(string)),
				d["dayId"].(string),
				d["patientId"].(string))
		},
		func(v interface{}) map[string]interface{} {
			t := v.(events.SlotBooked)
			return map[string]interface{}{
				"slotId":    t.SlotId,
				"dayId":     t.DayId,
				"patientId": t.PatientId,
			}
		})

	tm.MapEvent(infrastructure.GetValueType(events.SlotBookingCancelled{}), Prefix+"-slot-booking-cancelled",
		func(d map[string]interface{}) interface{} {
			return events.NewSlotBookingCancelled(
				uuid.MustParse(d["slotId"].(string)),
				d["dayId"].(string),
				d["reason"].(string))
		},
		func(v interface{}) map[string]interface{} {
			t := v.(events.SlotBookingCancelled)
			return map[string]interface{}{
				"slotId": t.SlotId,
				"dayId":  t.DayId,
				"reason": t.Reason,
			}
		})

	tm.MapEvent(infrastructure.GetValueType(events.SlotScheduleCancelled{}), Prefix+"-slot-schedule-cancelled",
		func(d map[string]interface{}) interface{} {
			return events.NewSlotScheduleCancelled(
				uuid.MustParse(d["slotId"].(string)),
				d["dayId"].(string))
		},
		func(v interface{}) map[string]interface{} {
			t := v.(events.SlotScheduleCancelled)
			return map[string]interface{}{
				"slotId": t.SlotId,
				"dayId":  t.DayId,
			}
		})

	tm.MapEvent(infrastructure.GetValueType(events.DayScheduleCancelled{}), Prefix+"-day-schedule-cancelled",
		func(d map[string]interface{}) interface{} {
			return events.NewDayScheduleCancelled(
				d["dayId"].(string))
		},
		func(v interface{}) map[string]interface{} {
			t := v.(events.DayScheduleCancelled)
			return map[string]interface{}{
				"dayId": t.DayId,
			}
		})

	tm.MapEvent(infrastructure.GetValueType(events.DayScheduleArchived{}), Prefix+"-day-schedule-archived",
		func(d map[string]interface{}) interface{} {
			return events.NewDayScheduleArchived(
				d["dayId"].(string))
		},
		func(v interface{}) map[string]interface{} {
			t := v.(events.DayScheduleArchived)
			return map[string]interface{}{
				"dayId": t.DayId,
			}
		})

	tm.MapEvent(infrastructure.GetValueType(events.CalendarDayStarted{}), Prefix+"-calendar-day-started",
		func(d map[string]interface{}) interface{} {
			return events.NewCalendarDayStarted(
				mustParseDate(d["date"].(string)))
		},
		func(v interface{}) map[string]interface{} {
			t := v.(events.CalendarDayStarted)
			return map[string]interface{}{
				"date": t.Date,
			}
		})

	registerType := func(t interface{}, typeName string) {
		tm.RegisterType(infrastructure.GetValueType(t), typeName, func() interface{} {
			return t
		})
	}

	// Commands
	registerType(commands.ArchiveDaySchedule{}, Prefix+"-archive-day-schedule")
	registerType(commands.BookSlot{}, Prefix+"-book-slot")
	registerType(commands.CancelDaySchedule{}, Prefix+"-cancel-day-schedule")
	registerType(commands.CancelSlotBooking{}, Prefix+"-cancel-slot-booking")
	registerType(commands.ScheduleDay{}, Prefix+"-schedule-day")
	registerType(commands.ScheduledSlot{}, Prefix+"-schedule-slot")

	// Snapshots
	registerType(DaySnapshot{}, "doctor-day-snapshot")
}
