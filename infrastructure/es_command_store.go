package infrastructure

import (
	"context"

	"github.com/EventStore/EventStore-Client-Go/esdb"
)

const (
	CommandsStreamName = "async_command_handler-day"
)

type CommandStore interface {
	Send(command interface{}, m CommandMetadata) error
	Start() error
}

type EsCommandStore struct {
	eventStore EventStore
	esdb       *esdb.Client
	serde      *EsEventSerde
	dispatcher *Dispatcher
}

func NewEsCommandStore(s EventStore, c *esdb.Client, e *EsEventSerde, d *Dispatcher) *EsCommandStore {
	return &EsCommandStore{
		eventStore: s,
		esdb:       c,
		serde:      e,
		dispatcher: d,
	}
}

func (c *EsCommandStore) Send(command interface{}, m CommandMetadata) error {
	return c.eventStore.AppendCommand(CommandsStreamName, command, m)
}

func (c *EsCommandStore) Start() error {
	options := esdb.SubscribeToStreamOptions{
		From: esdb.End{},
	}
	sub, err := c.esdb.SubscribeToStream(context.TODO(), c.eventStore.GetFullStreamName(CommandsStreamName), options)
	if err != nil {
		return err
	}

	go func() {
		for {
			s := sub.Recv()
			if s.EventAppeared != nil {
				cmd, m, err := c.serde.DeserializeCommand(s.EventAppeared)
				if err != nil {
					if cmd != nil {
						panic(err)
					} else {
						// ignore unknown event type
						continue
					}
				}

				c.dispatcher.Dispatch(cmd, m)
			}

			if s.SubscriptionDropped != nil {
				panic(s.SubscriptionDropped.Error)
			}
		}
	}()

	return nil
}
