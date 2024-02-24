package streamline

import (
	"context"
)

type EventsProcessor interface {
	ProcessEvents([]Event) error
	StreamTo(context.Context, Dispatcher) error
}

type Dispatcher interface {
	Dispatch(eventName string, payload []byte) error
}
