package projections

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/EventStore/training-introduction-go/infrastructure"
)

type SubscriptionManager struct {
	esdb            *esdb.Client
	checkpointStore *infrastructure.EsCheckpointStore
	serde           *infrastructure.EsEventSerde
	subscriptions   []Subscription
	streamName      string
	typesByName     map[string]reflect.Type
	isAllStream     bool
}

func NewSubscriptionManager(esdb *esdb.Client, c *infrastructure.EsCheckpointStore, s *infrastructure.EsEventSerde,
	streamName string, subs ...Subscription) *SubscriptionManager {
	return &SubscriptionManager{
		esdb:            esdb,
		subscriptions:   subs,
		checkpointStore: c,
		serde:           s,
		streamName:      streamName,
		isAllStream:     streamName == "$all",
	}
}

func (m SubscriptionManager) Start(ctx context.Context) error {
	position, err := m.checkpointStore.GetCheckpoint()
	if err != nil && !errors.Is(err, &infrastructure.CheckpointNotFoundError{}) {
		return err
	}

	var sub *esdb.Subscription
	if m.isAllStream {
		sub, err = m.esdb.SubscribeToAll(ctx, m.getAllStreamOptions(position))
	} else {
		sub, err = m.esdb.SubscribeToStream(ctx, m.streamName, m.getStreamOptions(position))
	}
	if err != nil {
		return err
	}

	go func() {
		for {
			s := sub.Recv()
			if s.EventAppeared != nil {
				if s.EventAppeared.Event == nil {
					continue
				}

				eventType := s.EventAppeared.Event.EventType
				if strings.HasPrefix(eventType, "$") || strings.Contains(eventType, "async_command_handler") {
					continue
				}

				event, metadata, err := m.serde.Deserialize(s.EventAppeared)
				if err != nil {
					if event != nil {
						panic(err)
					} else {
						// ignore unknown event type
						continue
					}
				}

				for _, s := range m.subscriptions {
					s.Project(event, *metadata)
				}

				m.storeCheckpoint(s)
			}

			if s.SubscriptionDropped != nil {
				panic(s.SubscriptionDropped.Error)
			}
		}
	}()

	return nil
}

func (m SubscriptionManager) getAllStreamOptions(position *uint64) esdb.SubscribeToAllOptions {
	options := esdb.SubscribeToAllOptions{}
	if position == nil {
		options.From = esdb.Start{}
	} else {
		options.From = esdb.Position{
			Commit:  *position,
			Prepare: *position,
		}
	}
	return options
}

func (m SubscriptionManager) getStreamOptions(position *uint64) esdb.SubscribeToStreamOptions {
	options := esdb.SubscribeToStreamOptions{}
	if position == nil {
		options.From = esdb.Start{}
	} else {
		options.From = esdb.Revision(*position)
	}
	return options
}

func (m SubscriptionManager) storeCheckpoint(s *esdb.SubscriptionEvent) {
	var checkpoint uint64
	if m.isAllStream {
		checkpoint = s.EventAppeared.OriginalEvent().Position.Commit
	} else {
		checkpoint = s.EventAppeared.Event.EventNumber
	}
	err := m.checkpointStore.StoreCheckpoint(checkpoint)
	if err != nil {
		panic(err)
	}
}
