package infrastructure

import (
	"reflect"
)

type EventHandler interface {
	CanHandle(reflect.Type) bool
	Handle(reflect.Type, interface{}, EventMetadata) error
	GetHandledTypes() []reflect.Type
}

type EventHandlerBase struct {
	EventHandler

	handlers     []EventHandlerEnvelope
	handledTypes map[reflect.Type]bool
	types        []reflect.Type
}

func NewEventHandler() EventHandlerBase {
	return EventHandlerBase{
		handledTypes: make(map[reflect.Type]bool, 0),
	}
}

func (p *EventHandlerBase) When(event interface{}, handler func(interface{}, EventMetadata)error) {
	t := GetValueType(event)
	p.handlers = append(p.handlers, NewEventHandlerEnvelope(t, handler))
	p.handledTypes[t] = true
}

func (p *EventHandlerBase) CanHandle(t reflect.Type) bool {
	_, exists := p.handledTypes[t]
	return exists
}

func (p *EventHandlerBase) Handle(t reflect.Type, event interface{}, m EventMetadata) error {
	for _, h := range p.handlers {
		if h.Type == t {
			err := h.Handler(event, m)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type EventHandlerEnvelope struct {
	Type    reflect.Type
	Handler func(interface{}, EventMetadata)error
}

func NewEventHandlerEnvelope(t reflect.Type, handler func(interface{}, EventMetadata)error) EventHandlerEnvelope {
	return EventHandlerEnvelope{
		Type:    t,
		Handler: handler,
	}
}
