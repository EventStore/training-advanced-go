package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/google/uuid"
)

const (
	ScheduleDayDateFormat = "2006-01-02"
	ScheduleDayTimeFormat = "15:04:05"
)

type ScheduleDayRequest struct {
	DoctorId uuid.UUID       `json:"doctorId"`
	Date     ScheduleDayDate `json:"date"`
	Slots    []SlotRequest   `json:"slots"`
}

type SlotRequest struct {
	Duration  SlotDuration `json:"duration"`
	StartTime string       `json:"startTime"`
}

func (r *ScheduleDayRequest) ToCommand() (commands.ScheduleDay, error) {
	slots := make([]commands.ScheduledSlot, 0)
	date := r.Date.Truncate(24 * time.Hour)

	for _, slot := range r.Slots {
		duration := slot.Duration
		startTime, err := time.Parse(ScheduleDayTimeFormat, slot.StartTime)
		slotStartTime := date.Add(
			time.Duration(startTime.Hour())*time.Hour +
				time.Duration(startTime.Minute())*time.Minute +
				time.Duration(startTime.Second())*time.Second)

		if err != nil {
			return commands.ScheduleDay{}, err
		}

		slots = append(slots, commands.NewScheduledSlot(slotStartTime, duration.ToDuration()))
	}

	return commands.NewScheduleDay(r.DoctorId, r.Date.Time, slots), nil
}

//
// Define custom type for date format
//

type ScheduleDayDate struct {
	time.Time
}

func (d *ScheduleDayDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse(ScheduleDayDateFormat, s)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

func (d ScheduleDayDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(ScheduleDayDateFormat))
}

//
// Define custom type for time format
//

type SlotDuration struct {
	time.Time
}

func (d *SlotDuration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse(ScheduleDayTimeFormat, s)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

func (d SlotDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		fmt.Sprintf("%02d:%02d:%02d", d.Hour(), d.Minute(), d.Second()))
}

func (d SlotDuration) ToDuration() time.Duration {
	return time.Duration(d.Hour())*time.Hour +
		time.Duration(d.Minute())*time.Minute +
		time.Duration(d.Second())*time.Second
}
