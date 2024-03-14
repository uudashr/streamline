package streamline

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
)

type eventRecord struct {
	name    string
	payload []byte
}

type InMemoryEventStream struct {
	events chan eventRecord
}

func NewInMemoryEventStream() *InMemoryEventStream {
	return &InMemoryEventStream{
		events: make(chan eventRecord, 100),
	}
}

func (es *InMemoryEventStream) Publish(ctx context.Context, event Event) error {
	eventType := reflect.TypeOf(event)
	eventName, ok := TagValue(eventType)
	if !ok {
		return errors.New("streamline: missing streamline tag")
	}

	// TODO: doesn't have to be json
	payload, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	es.events <- eventRecord{
		name:    eventName,
		payload: payload,
	}
	return nil
}

func (es *InMemoryEventStream) StreamTo(ctx context.Context, recv Receiver) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-es.events:
			if err := recv.Receive(event.name, event.payload); err != nil {
				return err
			}
		}
	}
}
