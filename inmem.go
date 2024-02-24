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

func (ep *InMemoryEventsProcessor) ProcessEvents(events []Event) error {
	for _, event := range events {
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
	}
	return nil
}

func (ep *InMemoryEventsProcessor) StreamTo(ctx context.Context, d Dispatcher) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-ep.events:
			if err := d.Dispatch(event.name, event.payload); err != nil {
				return err
			}
		}
	}
}
