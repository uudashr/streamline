package streamline

import (
	"reflect"

	"github.com/uudashr/rebound"
)

type Mux struct {
	rb *rebound.Rebound
}

func NewMux() *Mux {
	return &Mux{
		rb: &rebound.Rebound{},
	}
}

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

func (m *Mux) Dispatch(name string, payload []byte) error {
	return m.rb.Dispatch(name, payload)
}
