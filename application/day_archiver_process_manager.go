package application

import (
	"time"

	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/eventsourcing"
	"github.com/EventStore/training-introduction-go/infrastructure"
)

type DayArchiverProcessManager struct {
	infrastructure.EventHandlerBase
}

func NewDayArchiverProcessManager(s eventsourcing.ColdStorage, a readmodel.ArchivableDaysRepository,
	c infrastructure.CommandStore, es infrastructure.EventStore, archiveThreshold time.Duration) *DayArchiverProcessManager {
	h := infrastructure.NewEventHandler()

	h.When(events.DayScheduled{}, func(e interface{}, _ infrastructure.EventMetadata) error {
		d := e.(events.DayScheduled)
		a.Add(readmodel.NewArchivableDay(d.DayId, d.Date))
		return nil
	})

	h.When(events.CalendarDayStarted{}, func(e interface{}, m infrastructure.EventMetadata) error {
		d := e.(events.CalendarDayStarted)
		archivableDays, err := a.FindAll(d.Date.Add(archiveThreshold))
		if err != nil {
			return err
		}
		for _, a := range archivableDays {
			err := c.Send(commands.NewArchiveDaySchedule(a.Id), infrastructure.NewCommandMetadataFrom(m))
			if err != nil {
				return err
			}
		}
		return nil
	})

	h.When(events.DayScheduleArchived{}, func(e interface{}, _ infrastructure.EventMetadata) error {
		//	d := e.(events.DayScheduleArchived)

		//	streamName := eventsourcing.GetStreamNameWithId(&doctorday.Day{}, d.DayId)

		return nil
	})

	return &DayArchiverProcessManager{h}
}
