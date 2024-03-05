package streamline

import (
	"context"
)

type EventStream interface {
	PublishEvents(context.Context, []Event) error
	StreamTo(context.Context, Dispatcher) error
}

type Dispatcher interface {
	Dispatch(eventName string, payload []byte) error
}
