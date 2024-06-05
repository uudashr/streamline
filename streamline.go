package streamline

import (
	"errors"
	"strings"

	"github.com/uudashr/eventually"
	"github.com/uudashr/rebound"
)

type Event = eventually.Event
type EventHandler = rebound.EventHandler

type Receiver interface {
	Receive(eventName string, payload []byte) error
}

type ReceiverFunc func(eventName string, payload []byte) error

func (fn ReceiverFunc) Receive(eventName string, payload []byte) error {
	return fn(eventName, payload)
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
