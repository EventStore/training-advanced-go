package doctorday

import "github.com/EventStore/training-introduction-go/infrastructure"

type EventStoreDayRepository struct {
	aggregateStore infrastructure.AggregateStore
}

func NewEventStoreDayRepository(store infrastructure.AggregateStore) *EventStoreDayRepository {
	return &EventStoreDayRepository{
		aggregateStore: store,
	}
}

func (r *EventStoreDayRepository) Save(day *Day, m infrastructure.CommandMetadata) {
	r.aggregateStore.Save(day, m)
}

func (r *EventStoreDayRepository) Get(id DayID) (*Day, error) {
	day := NewDay()
	err := r.aggregateStore.Load(id.Value, day)
	if err != nil {
		return nil, err
	}

	return day, nil
}
