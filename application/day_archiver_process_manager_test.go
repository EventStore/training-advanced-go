package application

import (
	"fmt"
	"testing"
	"time"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/EventStore/training-introduction-go/domain/doctorday"
	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/EventStore/training-introduction-go/domain/doctorday/events"
	"github.com/EventStore/training-introduction-go/eventsourcing"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/EventStore/training-introduction-go/infrastructure/inmemory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDayArchiver(t *testing.T) {
	client, err := createESDBClient()
	assert.NoError(t, err)

	typeMapper := eventsourcing.NewTypeMapper()
	esSerde := infrastructure.NewEsEventSerde(typeMapper)
	tenantPrefix := fmt.Sprintf("day_archiver_tests_%s_", uuid.NewString())
	eventStore := infrastructure.NewEsEventStore(client, tenantPrefix, esSerde)
	cmdStore := infrastructure.NewEsCommandStore(eventStore, nil, nil, nil)
	doctorday.RegisterTypes(typeMapper)

	p := &DayArchiverTests{
		HandlerTests: infrastructure.NewHandlerTests(t),

		dayId:       "dayId",
		patientId:   "patientId",
		reason:      "Some cancellation reason",
		slotId:      uuid.New(),
		now:         time.Now().Truncate(time.Second),
		tenMinutes:  10 * time.Minute,
		eventStore:  eventStore,
		coldStorage: inmemory.NewColdStorage(),
	}

	p.SetHandlerFactory(func() infrastructure.EventHandler {
		r := inmemory.NewArchivableDaysRepository()
		return NewDayArchiverProcessManager(p.coldStorage, r, cmdStore, eventStore, time.Hour*-24)
	})

	t.Run("ShouldArchiveAllEventsAndTruncateAllExceptLastOne", p.ShouldArchiveAllEventsAndTruncateAllExceptLastOne)
	t.Run("ShouldSendArchiveCommandForAllSlotsCompleted180DaysAgo", p.ShouldSendArchiveCommandForAllSlotsCompleted180DaysAgo)
}

func (a *DayArchiverTests) ShouldArchiveAllEventsAndTruncateAllExceptLastOne(t *testing.T) {
	dayId := uuid.NewString()
	scheduled := events.NewSlotScheduled(uuid.New(), dayId, a.now, a.tenMinutes)
	slotBooked := events.NewSlotBooked(scheduled.SlotId, dayId, "PatientId")
	dayArchived := events.NewDayScheduleArchived(dayId)
	metadata := infrastructure.NewCommandMetadata(uuid.New(), uuid.New())

	events := []interface{}{scheduled, slotBooked, dayArchived}

	streamName := eventsourcing.GetStreamNameWithId(&doctorday.Day{}, dayId)
	err := a.eventStore.AppendEventsToAny(streamName, metadata, events...)
	assert.NoError(t, err)

	a.Given(dayArchived)
	a.Then(events, a.coldStorage.Events)

	loadedEvents, err := a.eventStore.LoadEventsFromStart(streamName)
	assert.NoError(t, err)
	assert.Len(t, loadedEvents, 1)

	a.Then(dayArchived, loadedEvents[0])
}

func (a *DayArchiverTests) ShouldSendArchiveCommandForAllSlotsCompleted180DaysAgo(t *testing.T) {
	dayId := uuid.NewString()
	date := time.Now().AddDate(0, 0, -180)
	dayScheduled := events.NewDayScheduled(dayId, uuid.New(), date)
	calenderDayStarted := events.NewCalendarDayStarted(a.now)

	a.Given(dayScheduled, calenderDayStarted)

	cmds, err := a.eventStore.LoadCommand("async_command_handler-day")
	assert.NoError(t, err)
	assert.Len(t, cmds, 1)

	a.Then(
		commands.NewArchiveDaySchedule(dayId),
		cmds[0].Command)
}

func createESDBClient() (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString("esdb://localhost:2113?tls=false")
	if err != nil {
		return nil, err
	}

	db, err := esdb.NewClient(settings)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type DayArchiverTests struct {
	infrastructure.HandlerTests

	eventStore  *infrastructure.EsEventStore
	coldStorage *inmemory.ColdStorage

	dayId      string
	patientId  string
	reason     string
	slotId     uuid.UUID
	now        time.Time
	tenMinutes time.Duration
}
