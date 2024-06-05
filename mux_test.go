package streamline_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/uudashr/streamline"
)

func ExampleMux() {
	mux := streamline.NewMux()
	type OrderCompleted struct {
		OrderID string `json:"orderId" streamline:"order.completed"` // defines the event name also the field is marked as object/aggregate id
	}

	// It will react to the OrderCompleted event. This can be taken from the messaging service.
	mux.React(func(event OrderCompleted) error {
		fmt.Println("Got OrderCompleted event:", event.OrderID)
		return nil
	})

	mux.Dispatch("order.completed", []byte(`{"orderId":"123"}`))
	// Output: Got OrderCompleted event: 123
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
