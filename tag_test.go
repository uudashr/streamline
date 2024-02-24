package streamline_test

import (
	"reflect"
	"testing"

	"github.com/uudashr/streamline"
)

func TestTagValue(t *testing.T) {
	type OrderCompleted struct {
		OrderID string `streamline:"order.completed"`
	}

	event := OrderCompleted{
		OrderID: "order-123",
	}

	tagVal, ok := streamline.TagValue(reflect.TypeOf(event))
	if !ok {
		t.Error("expected ok to be true")
	}

	if got, want := tagVal, "order.completed"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTagFieldValue(t *testing.T) {
	type OrderCompleted struct {
		OrderID string `streamline:"order.completed"`
	}

	event := OrderCompleted{
		OrderID: "order-123",
	}

	fieldVal, ok := streamline.TagFieldValue(event)
	if !ok {
		t.Error("expected ok to be true")
	}

	if got, want := fieldVal, "order-123"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
