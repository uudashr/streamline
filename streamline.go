package streamline

import (
	"context"
	"errors"
	"reflect"

	"github.com/uudashr/eventually"
	"github.com/uudashr/rebound"
)

type Event = eventually.Event
type EventHandler = rebound.EventHandler

type Streamline struct {
	rb              *rebound.Rebound
	eventsProcessor EventStream
}

func NewStreamline(rb *rebound.Rebound, eventsProcessor EventStream) (*Streamline, error) {
	if rb == nil {
		return nil, errors.New("streamline: nil rb")
	}

	if eventsProcessor == nil {
		return nil, errors.New("streamline: nil eventsProcessor")
	}

	return &Streamline{
		rb:              rb,
		eventsProcessor: eventsProcessor,
	}, nil
}

func (stln *Streamline) React(fn EventHandler) {
	fnt := reflect.TypeOf(fn)
	inType := fnt.In(0)
	if inType.Kind() != reflect.Struct {
		panic("streamline: fn EventHandler argument must be a struct")
	}

	eventName, ok := TagValue(inType)
	if !ok {
		panic("streamline: missing streamline tag")
	}

	stln.rb.ReactTo(eventName, fn)
}

func (stln *Streamline) ProcessEvents(ctx context.Context, events []Event) error {
	// Options:
	// 1. Publish to a message service
	// 2. Process in-memory
	// 3. Persist to a database (outbox table)
	return stln.eventsProcessor.PublishEvents(ctx, events)
}

// Deprecated: use rebound.Rebound instance directly
func (stln *Streamline) Dispatch(eventName string, payload []byte) error {
	return stln.rb.Dispatch(eventName, payload)
}
