package infrastructure

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type AggregateTests struct {
	dispatcher  Dispatcher
	store       *FakeAggregateStore
	latestError error
}

func NewAggregateTests(store *FakeAggregateStore) AggregateTests {
	return AggregateTests{
		store: store,
	}
}

func (t *AggregateTests) RegisterHandlers(handlers CommandHandler) {
	commandHandlerMap := NewCommandHandlerMap(handlers)
	t.dispatcher = NewDispatcher(commandHandlerMap)
}

func (t *AggregateTests) Given(events... interface{}) {
	t.store.SetInitialEvents(events)
}

func (t *AggregateTests) When(command interface{}) {
	t.latestError = t.dispatcher.Dispatch(command, NewCommandMetadata(uuid.New(), uuid.New()))
}

func (t *AggregateTests) Then(then func([]interface{}, error)) {
	then(t.store.GetStoredChanges(), t.latestError)
}

func (t *AggregateTests) ThenExpectError(tt *testing.T, expected error) {
	assert.Error(tt, t.latestError)
	assert.Equal(tt, expected, t.latestError)
}

func (t *AggregateTests) ThenExpectSingleChange(tt *testing.T, expected interface{}) {
	t.ThenExpectChange(tt, 0, expected)
}

func (t *AggregateTests) ThenExpectChanges(tt *testing.T, expectedChanges []interface{}) {
	for i, expected := range expectedChanges {
		t.ThenExpectChange(tt, i, expected)
	}
}

func (t *AggregateTests) ThenExpectChange(tt *testing.T, idx int, expected interface{}) {
	changes := t.store.GetStoredChanges()
	assert.NoError(tt, t.latestError)
	assert.IsType(tt, expected, changes[idx])
	assert.Equal(tt, expected, changes[idx])
}
