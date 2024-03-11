package streamline

import (
	"context"
)

type EventStream interface {
	Publish(context.Context, Event) error
	StreamTo(context.Context, Receiver) error
}

type Receiver interface {
	Receive(eventName string, payload []byte) error
}

type ReceiverFunc func(eventName string, payload []byte) error

func (fn ReceiverFunc) Receive(eventName string, payload []byte) error {
	return fn(eventName, payload)
}
