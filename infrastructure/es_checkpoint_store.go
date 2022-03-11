package infrastructure

import (
	"context"
	"encoding/json"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gofrs/uuid"
)

const (
	CheckpointStreamPrefix = "checkpoint-"
)

type EsCheckpointStore struct {
	esdb       *esdb.Client
	serde      *EsEventSerde
	streamName string
}

func NewEsCheckpointStore(esdb *esdb.Client, subscriptionName string, serde *EsEventSerde) *EsCheckpointStore {
	return &EsCheckpointStore{
		esdb:       esdb,
		serde:      serde,
		streamName: CheckpointStreamPrefix + subscriptionName,
	}
}

func (s *EsCheckpointStore) GetCheckpoint() (*uint64, error) {
	options := esdb.ReadStreamOptions{
		From:      esdb.End{},
		Direction: esdb.Backwards,
	}
	result, err := s.esdb.ReadStream(context.TODO(), s.streamName, options, 1)
	if err != nil {
		if esdbError, ok := esdb.FromError(err); !ok {
			if esdbError.Code() == esdb.ErrorResourceNotFound {
				return nil, &CheckpointNotFoundError{}
			}
		}
		return nil, err
	}

	e, err := result.Recv()
	if err != nil {
		return nil, err
	}

	if e.Event != nil {
		c := Checkpoint{}
		err = json.Unmarshal(e.Event.Data, &c)
		if err != nil {
			return nil, err
		}

		return &c.Position, nil
	}

	return nil, &CheckpointNotFoundError{}
}

func (s *EsCheckpointStore) StoreCheckpoint(position uint64) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	checkpointData, err := json.Marshal(Checkpoint{Position: position})
	if err != nil {
		return err
	}

	eventData := esdb.EventData{
		EventID:     id,
		ContentType: esdb.JsonContentType,
		EventType:   "$checkpoint",
		Data:        checkpointData,
	}

	options := esdb.AppendToStreamOptions{
		ExpectedRevision: esdb.Any{},
	}

	_, err = s.esdb.AppendToStream(context.TODO(), s.streamName, options, eventData)
	return err
}

type Checkpoint struct {
	Position uint64 `json:"position"`
}

type CheckpointNotFoundError struct{}

func (c *CheckpointNotFoundError) Error() string {
	return "checkpoint not found"
}
