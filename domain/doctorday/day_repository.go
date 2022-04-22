package doctorday

import (
	"github.com/EventStore/training-introduction-go/infrastructure"
)

type DayRepository interface {
	Save(day *Day, metadata infrastructure.CommandMetadata)
	Get(id DayID) (*Day, error)
}
