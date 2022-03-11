package infrastructure

import (
	"fmt"
)

type Dispatcher struct {
	commandHandlerMap CommandHandlerMap
}

func (d Dispatcher) Dispatch(command interface{}, metadata CommandMetadata) error {
	handler, err := d.commandHandlerMap.Get(GetValueType(command))
	if err != nil {
		return fmt.Errorf("no handler registered")
	}

	return handler(command, metadata)
}

func NewDispatcher(commandHandlerMap CommandHandlerMap) Dispatcher {
	return Dispatcher{
		commandHandlerMap: commandHandlerMap,
	}
}
