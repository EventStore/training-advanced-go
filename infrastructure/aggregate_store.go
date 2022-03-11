package infrastructure

import (
	"fmt"

	"github.com/EventStore/training-introduction-go/eventsourcing"
)

type AggregateStore interface {
	Save(a eventsourcing.AggregateRoot, m CommandMetadata) error
	Load(id string, a eventsourcing.AggregateRoot) error
}

type AggregateNotFoundError struct{}

func (e AggregateNotFoundError) Error() string {
	return fmt.Sprintf("aggregate not found error")
}
