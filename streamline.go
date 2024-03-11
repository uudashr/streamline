package streamline

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/uudashr/eventually"
	"github.com/uudashr/rebound"
)

type Event = eventually.Event
type EventHandler = rebound.EventHandler

type Streamline struct {
	rb     *rebound.Rebound
	stream EventStream
}

func NewStreamline(rb *rebound.Rebound, stream EventStream) (*Streamline, error) {
	if rb == nil {
		return nil, errors.New("streamline: nil rb")
	}

	if stream == nil {
		return nil, errors.New("streamline: nil stream")
	}

	return &Streamline{
		rb:     rb,
		stream: stream,
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

func (stln *Streamline) Publish(ctx context.Context, event Event) error {
	// Options:
	// 1. Publish to a message service
	// 2. Process in-memory
	// 3. Persist to a database (outbox table)
	return stln.stream.Publish(ctx, event)
}

type Meta struct {
	ObjectName string
	ObjectID   string
	EventName  string
}

func MetaOf(event Event) (Meta, error) {
	tag, val, ok := Tag(event)
	if !ok {
		return Meta{}, errors.New("streamline: no tag found")
	}

	vals := strings.Split(tag, ".")
	if len(vals) != 2 {
		return Meta{}, errors.New("streamline: invalid tag")
	}

	return Meta{
		ObjectName: vals[0],
		ObjectID:   val,
		EventName:  vals[1],
	}, nil
}
