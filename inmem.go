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

type InMemoryEventsProcessor struct {
	events chan eventRecord
}

func NewInMemoryEventsProcessor() *InMemoryEventsProcessor {
	return &InMemoryEventsProcessor{
		events: make(chan eventRecord, 100),
	}
}

func (ep *InMemoryEventsProcessor) Publish(ctx context.Context, event Event) error {
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

	ep.events <- eventRecord{
		name:    eventName,
		payload: payload,
	}
	return nil
}

func (ep *InMemoryEventsProcessor) StreamTo(ctx context.Context, recv Receiver) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-ep.events:
			if err := recv.Receive(event.name, event.payload); err != nil {
				return err
			}
		}
	}
}
