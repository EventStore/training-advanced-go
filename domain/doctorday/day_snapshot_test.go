package doctorday

import (
	"testing"
	"time"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/EventStore/training-introduction-go/domain/doctorday/commands"
	"github.com/EventStore/training-introduction-go/eventsourcing"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWriteSnapshotIfThresholdReached(t *testing.T) {
	client, err := createESDBClient()
	assert.NoError(t, err)

	defer client.Close()

	typeMapper := eventsourcing.NewTypeMapper()
	esSerde := infrastructure.NewEsEventSerde(typeMapper)
	store := infrastructure.NewEsEventStore(client, "snapshot_test-", esSerde)
	aggregateStore := infrastructure.NewEsAggregateStore(store, 5)
	RegisterTypes(typeMapper)

	now := time.Now()
	tenMinutes := time.Duration(10) * time.Minute
	slots := []commands.ScheduledSlot{
		commands.NewScheduledSlot(now, tenMinutes),
		commands.NewScheduledSlot(now.Add(tenMinutes), tenMinutes),
		commands.NewScheduledSlot(now.Add(tenMinutes*2), tenMinutes),
		commands.NewScheduledSlot(now.Add(tenMinutes*3), tenMinutes),
		commands.NewScheduledSlot(now.Add(tenMinutes*4), tenMinutes),
	}

	aggregate := NewDay()
	aggregate.ScheduleDay(NewDoctorID(uuid.New()), now, slots)

	err = aggregateStore.Save(aggregate, infrastructure.NewCommandMetadata(uuid.New(), uuid.New()))
	assert.NoError(t, err)

	streamName := eventsourcing.GetStreamName(aggregate)
	s, m, err := store.LoadSnapshot(streamName)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, m)
}

func TestReadSnapshotWhenLoadingAggregate(t *testing.T) {
	client, err := createESDBClient()
	assert.NoError(t, err)

	defer client.Close()

	typeMapper := eventsourcing.NewTypeMapper()
	esSerde := infrastructure.NewEsEventSerde(typeMapper)
	store := infrastructure.NewEsEventStore(client, "snapshot_test-", esSerde)
	aggregateStore := infrastructure.NewEsAggregateStore(store, 5)
	RegisterTypes(typeMapper)

	now := time.Now()
	tenMinutes := time.Duration(10) * time.Minute
	slots := []commands.ScheduledSlot{
		commands.NewScheduledSlot(now, tenMinutes),
		commands.NewScheduledSlot(now.Add(tenMinutes), tenMinutes),
		commands.NewScheduledSlot(now.Add(tenMinutes*2), tenMinutes),
		commands.NewScheduledSlot(now.Add(tenMinutes*3), tenMinutes),
		commands.NewScheduledSlot(now.Add(tenMinutes*4), tenMinutes),
	}

	aggregate := NewDay()
	aggregate.ScheduleDay(NewDoctorID(uuid.New()), now, slots)
	aggregateChangeCount := len(aggregate.GetChanges())

	err = aggregateStore.Save(aggregate, infrastructure.NewCommandMetadata(uuid.New(), uuid.New()))
	assert.NoError(t, err)

	streamName := eventsourcing.GetStreamName(aggregate)
	err = store.TruncateStream(streamName, uint64(aggregateChangeCount))
	assert.NoError(t, err)

	reloadedAggregate := NewDay()
	err = aggregateStore.Load(aggregate.Id, reloadedAggregate)
	assert.NoError(t, err)
	assert.Equal(t, 5, reloadedAggregate.GetVersion())
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
