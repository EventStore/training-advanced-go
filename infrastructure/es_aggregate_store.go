package infrastructure

import (
	"github.com/EventStore/training-introduction-go/eventsourcing"
)

type EsAggregateStore struct {
	AggregateStore

	store             EventStore
	snapshotThreshold int
}

func NewEsAggregateStore(store EventStore, snapshotThreshold int) *EsAggregateStore {
	return &EsAggregateStore{
		store:             store,
		snapshotThreshold: snapshotThreshold,
	}
}

func (s *EsAggregateStore) Save(a eventsourcing.AggregateRoot, m CommandMetadata) error {
	changes := a.GetChanges()
	streamName := eventsourcing.GetStreamName(a)
	err := s.store.AppendEvents(streamName, a.GetVersion(), m, changes...)
	if err != nil {
		return err
	}

	if sa, ok := a.(eventsourcing.AggregateRootSnapshot); ok {
		newVersion := a.GetVersion() + len(changes)
		if (newVersion+1)-sa.GetSnapshotVersion() >= s.snapshotThreshold {
			err = s.store.AppendSnapshot(streamName, newVersion, sa.GetSnapshot())
			if err != nil {
				return err
			}
		}
	}

	a.ClearChanges()
	return nil
}

func (s *EsAggregateStore) Load(aggregateId string, a eventsourcing.AggregateRoot) error {
	version := -1
	streamName := eventsourcing.GetStreamNameWithId(a, aggregateId)

	if sa, ok := a.(eventsourcing.AggregateRootSnapshot); ok {
		sn, md, _ := s.store.LoadSnapshot(streamName)
		if sn != nil && md != nil {
			sa.LoadSnapshot(sn, md.Version)
			version = md.Version + 1
		}
	}

	events, err := s.store.LoadEvents(streamName, version)
	if err != nil {
		return err
	}

	a.Load(events)
	a.ClearChanges()
	return nil
}
