package streamline

import (
	"errors"
	"strings"

	"github.com/uudashr/eventually"
	"github.com/uudashr/rebound"
)

// Event is a type of event. It can be struct of anything represent the event.
//
// Event should contains "streamline" tag in the form of fully qualified event name in the form of `objectName.eventName`.
// The tag must be placed on the field that represent the object ID.
//
// Example:
//
//	type OrderPlaced struct {
//	  OrderID string `streamline:"order.OrderPlaced"`
//	  CustomerID string
//	}
//
// See more: [eventually.Event].
type Event = eventually.Event

// Eventhandler is a function type that handles an event.
//
// It is a function that has one argument of event type and returns an error.
//
// See [rebound.EventHandler].
type EventHandler = rebound.EventHandler

// Receiver receives an event in a bytes payload form identified by the event name.
type Receiver interface {
	Receive(eventName string, payload []byte) error
}

// ReceiverFunc is a function type that implements [Receiver] interface.
type ReceiverFunc func(eventName string, payload []byte) error

// Receive implement the [Receiver] interface.
func (fn ReceiverFunc) Receive(eventName string, payload []byte) error {
	return fn(eventName, payload)
}

// Meta is the meta data of the event.
type Meta struct {
	ObjectName string
	ObjectID   string
	EventName  string
}

// MetaOf returns the meta data of the event.
// The event must have a tag in the form of `objectName.eventName`.
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
