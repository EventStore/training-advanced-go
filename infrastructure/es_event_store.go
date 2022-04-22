package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/EventStore/training-introduction-go/eventsourcing"
)

//todo move to /event-sourcing

type CommandEnvelope struct {
	Command  interface{}
	Metadata CommandMetadata
}

func NewCommandEnvelope(command interface{}, m CommandMetadata) CommandEnvelope {
	return CommandEnvelope{
		Command:  command,
		Metadata: m,
	}
}

type EventStore interface {
	GetFullStreamName(streamName string) string
	GetLastVersion(streamName string) (uint64, error)
	TruncateStream(streamName string, version uint64) error

	AppendEvents(streamName string, version int, m CommandMetadata, events ...interface{}) error
	AppendEventsToAny(streamName string, m CommandMetadata, events ...interface{}) error
	LoadEvents(streamName string, version int) ([]interface{}, error)
	LoadEventsFromStart(streamName string) ([]interface{}, error)

	AppendSnapshot(streamName string, version int, snapshot interface{}) error
	LoadSnapshot(streamName string) (interface{}, *eventsourcing.SnapshotMetadata, error)

	AppendCommand(streamName string, command interface{}, m CommandMetadata) error
	LoadCommand(streamName string) ([]CommandEnvelope, error)
}

type EsEventStore struct {
	EventStore

	esdb         *esdb.Client
	tenantPrefix string
	esSerde      *EsEventSerde
}

func NewEsEventStore(esdb *esdb.Client, tenantPrefix string, serder *EsEventSerde) *EsEventStore {
	return &EsEventStore{
		esdb:         esdb,
		tenantPrefix: tenantPrefix,
		esSerde:      serder,
	}
}

func (s *EsEventStore) GetFullStreamName(streamName string) string {
	return s.tenantPrefix + streamName
}

func (s *EsEventStore) GetLastVersion(streamName string) (uint64, error) {
	options := esdb.ReadStreamOptions{
		From:      esdb.End{},
		Direction: esdb.Backwards,
	}
	result, err := s.esdb.ReadStream(context.TODO(), s.GetFullStreamName(streamName), options, 1)
	if err != nil {
		return 0, err
	}

	e, err := result.Recv()
	if err != nil {
		return 0, err
	}

	if e.Event != nil {
		return e.Event.EventNumber, nil
	}

	return 0, fmt.Errorf("failed to retrieve last version")
}

func (s *EsEventStore) AppendEventsToAny(streamName string, m CommandMetadata, events ...interface{}) error {
	options := esdb.AppendToStreamOptions{
		ExpectedRevision: esdb.Any{},
	}
	return s.appendEvents(streamName, options, m, events...)
}

func (s *EsEventStore) AppendEvents(streamName string, version int, m CommandMetadata, events ...interface{}) error {
	options := esdb.AppendToStreamOptions{}
	if version == -1 {
		options.ExpectedRevision = esdb.NoStream{}
	} else {
		options.ExpectedRevision = esdb.Revision(uint64(version))
	}

	return s.appendEvents(streamName, options, m, events...)
}

func (s *EsEventStore) appendEvents(streamName string, o esdb.AppendToStreamOptions, m CommandMetadata, events ...interface{}) error {
	if events == nil || len(events) == 0 {
		return nil
	}

	var eventData []esdb.EventData
	for _, e := range events {
		ed, err := s.esSerde.Serialize(e, NewEventMetadataFrom(m))
		if err != nil {
			return err
		}
		eventData = append(eventData, ed)
	}

	_, err := s.esdb.AppendToStream(context.TODO(), s.GetFullStreamName(streamName), o, eventData...)
	return err
}

func (s *EsEventStore) LoadEventsFromStart(streamName string) ([]interface{}, error) {
	options := esdb.ReadStreamOptions{
		From:      esdb.Start{},
		Direction: esdb.Forwards,
	}

	return s.loadEvents(streamName, options)
}

func (s *EsEventStore) LoadEvents(streamName string, version int) ([]interface{}, error) {
	options := esdb.ReadStreamOptions{
		Direction: esdb.Forwards,
	}

	if version == -1 {
		options.From = esdb.Start{}
	} else {
		options.From = esdb.Revision(uint64(version))
	}

	return s.loadEvents(streamName, options)
}

func (s *EsEventStore) loadEvents(streamName string, o esdb.ReadStreamOptions) ([]interface{}, error) {
	events := make([]interface{}, 0)
	stream, err := s.esdb.ReadStream(context.TODO(), s.GetFullStreamName(streamName), o, math.MaxInt64)
	if err != nil {
		if esdbError, ok := esdb.FromError(err); !ok {
			if esdbError.Code() == esdb.ErrorResourceNotFound {
				return events, nil
			} else if errors.Is(err, io.EOF) {
				return events, nil
			}
		}
		return nil, err
	}

	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		e, _, err := s.esSerde.Deserialize(event)
		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

func (s *EsEventStore) AppendSnapshot(streamName string, version int, snapshot interface{}) error {
	eventData, err := s.esSerde.SerializeSnapshot(snapshot, eventsourcing.NewSnapshotMetadata(version))
	if err != nil {
		return err
	}

	options := esdb.AppendToStreamOptions{ExpectedRevision: esdb.Any{}}
	snapshotName := s.GetFullStreamName("snapshot-" + streamName)
	_, err = s.esdb.AppendToStream(context.Background(), snapshotName, options, eventData)
	return err
}

func (s *EsEventStore) LoadSnapshot(streamName string) (interface{}, *eventsourcing.SnapshotMetadata, error) {
	options := esdb.ReadStreamOptions{
		From:      esdb.End{},
		Direction: esdb.Backwards,
	}
	snapshotName := s.GetFullStreamName("snapshot-" + streamName)
	stream, err := s.esdb.ReadStream(context.Background(), snapshotName, options, 1)
	if err != nil {
		return nil, nil, err
	}

	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil, nil, fmt.Errorf("unexpected end of stream")
		}

		if err != nil {
			return nil, nil, err
		}

		return s.esSerde.DeserializeSnapshot(event)
	}

	return nil, nil, fmt.Errorf("failed to load snapshot")
}

func (s *EsEventStore) TruncateStream(streamName string, beforeVersion uint64) error {
	options := esdb.AppendToStreamOptions{ExpectedRevision: esdb.Any{}}
	metadata := esdb.StreamMetadata{}
	metadata.SetTruncateBefore(beforeVersion)
	_, err := s.esdb.SetStreamMetadata(context.Background(), s.GetFullStreamName(streamName), options, metadata)
	return err
}

func (s *EsEventStore) AppendCommand(streamName string, command interface{}, m CommandMetadata) error {
	eventData, err := s.esSerde.SerializeCommand(command, m)
	if err != nil {
		return err
	}

	options := esdb.AppendToStreamOptions{ExpectedRevision: esdb.Any{}}
	_, err = s.esdb.AppendToStream(context.TODO(), s.GetFullStreamName(streamName), options, eventData)
	return err
}

func (s *EsEventStore) LoadCommand(streamName string) ([]CommandEnvelope, error) {
	options := esdb.ReadStreamOptions{
		From:      esdb.Start{},
		Direction: esdb.Forwards,
	}
	stream, err := s.esdb.ReadStream(context.TODO(), s.GetFullStreamName(streamName), options, math.MaxUint64)
	if err != nil {
		return nil, err
	}

	cmdEnvelopes := make([]CommandEnvelope, 0)
	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		cmd, metadata, err := s.esSerde.DeserializeCommand(event)
		if err != nil {
			return nil, err
		}

		cmdEnvelopes = append(cmdEnvelopes, NewCommandEnvelope(cmd, metadata))
	}

	return cmdEnvelopes, nil
}
