package streamline_test

import (
	"reflect"
	"testing"

	"github.com/uudashr/rebound"
	"github.com/uudashr/streamline"
)

func ExampleStreamline() {
	var (
		rb     *rebound.Rebound
		stream streamline.EventStream
	)

	stln, _ := streamline.NewStreamline(rb, stream)

	type OrderCompleted struct {
		OrderID string `streamline:"order.completed"` // defines the event name also the field is marked as object/aggregate id
	}

	// It will react to the OrderCompleted event. This can be taken from the messaging service.
	stln.React(func(event OrderCompleted) error {
		return nil
	})

	rb.Dispatch("order.completed", []byte(`{"order_id":"123"}`))
}

func TestReflect(t *testing.T) {
	type OrderCompleted struct {
		OrderID string `streamline:"order.completed"`
	}

	fn := func(event OrderCompleted) error {
		return nil
	}

	fnt := reflect.TypeOf(fn)
	if got, want := fnt.In(0).Name(), "OrderCompleted"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	inType := fnt.In(0)
	tagVal := inType.Field(0).Tag.Get("streamline")
	if got, want := tagVal, "order.completed"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}
