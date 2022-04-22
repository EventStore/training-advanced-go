package infrastructure

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type HandlerTests struct {
	T *testing.T

	latestError    error
	handlerFactory func() EventHandler

	EnableAtLeastOnceMonkey  bool
	EnableAtLeastOnceGorilla bool
}

func NewHandlerTests(t *testing.T) HandlerTests {
	return HandlerTests{T: t}
}

func (p *HandlerTests) SetHandlerFactory(f func() EventHandler) {
	p.handlerFactory = f
}

func (p *HandlerTests) Given(events ...interface{}) {
	assert.NotNil(p.T, p.handlerFactory)

	correlationId := uuid.New()
	causationId := uuid.New()

	eventHandler := p.handlerFactory()
	for i, e := range events {
		m := NewEventMetadata(correlationId, causationId, i)

		err := eventHandler.Handle(GetValueType(e), e, m)
		if err != nil {
			p.latestError = err
			return
		}

		if p.EnableAtLeastOnceMonkey {
			err = eventHandler.Handle(GetValueType(e), e, m)
			if err != nil {
				p.latestError = err
				return
			}
		}
	}

	if p.EnableAtLeastOnceGorilla {
		for _, e := range events[:len(events)-1] {

			m := NewEventMetadata(correlationId, causationId, 7)

			err := eventHandler.Handle(GetValueType(e), e, m)
			if err != nil {
				p.latestError = err
				return
			}

			if p.EnableAtLeastOnceMonkey {
				err = eventHandler.Handle(GetValueType(e), e, m)
				if err != nil {
					p.latestError = err
					return
				}
			}
		}
	}
}

func (p *HandlerTests) Then(expected, actual interface{}) {
	assert.NoError(p.T, p.latestError)
	assert.Equal(p.T, expected, actual)
}
