package streamline

import (
	"errors"
	"reflect"

	"github.com/uudashr/rebound"
)

// NoHandlerError is an error when no handler found for the event.
//
// See [rebound.NoHandlerError].
type NoHandlerError = rebound.NoHandlerError

// Mux is multiplexer that dispatches the event to the registered handler.
type Mux struct {
	rb              *rebound.Rebound
	IgnoreNoHandler bool
}

// NewMux creates a new Mux.
func NewMux() *Mux {
	return &Mux{
		rb: &rebound.Rebound{},
	}
}

// React to an event which defined by the fn handler.
func (m *Mux) React(fn EventHandler) {
	fnt := reflect.TypeOf(fn)

	inType := fnt.In(0)
	if inType.Kind() != reflect.Struct {
		panic("streamline: fn EventHandler argument must be a struct")
	}

	eventName, ok := TagValue(inType)
	if !ok {
		panic("streamline: missing streamline tag")
	}

	m.rb.ReactTo(eventName, fn)
}

// Dispatch the event payload.
func (m *Mux) Dispatch(name string, payload []byte) error {
	err := m.rb.Dispatch(name, payload)

	var noHandlerErr NoHandlerError
	if m.IgnoreNoHandler && errors.As(err, &noHandlerErr) {
		return nil
	}

	return err
}
