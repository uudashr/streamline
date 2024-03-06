package streamline

import (
	"context"
)

type EventStream interface {
	Publish(context.Context, Event) error
	StreamTo(context.Context, Dispatcher) error
}

type Dispatcher interface {
	Dispatch(eventName string, payload []byte) error
}
