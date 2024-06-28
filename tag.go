package streamline

import (
	"reflect"
)

const (
	tagName = "streamline"
)

// TagValue returns the value of the streamline tag from the struct type.
func TagValue(inType reflect.Type) (value string, ok bool) {
	for i := range inType.NumField() {
		val := inType.Field(i).Tag.Get(tagName)
		if val != "" {
			return val, true
		}
	}

	return "", false
}

// TagFieldValue returns the value of the tag from the event.
func TagFieldValue(event Event) (value string, ok bool) {
	_, val, ok := Tag(event)

	return val, ok
}

// Tag returns the tag and the value of the tag from the event.
//
// Example:
//
//	type OrderCompleted struct {
//	  OrderID string `streamline:"order.completed"`
//	}
//
// Output:
//   - tag: "order.completed"
//   - value: "order-123"
//   - ok: true
func Tag(event Event) (tag string, value string, ok bool) {
	eventVal := reflect.ValueOf(event)
	for i := range eventVal.NumField() {
		tag, ok := eventVal.Type().Field(i).Tag.Lookup(tagName)
		if ok {
			return tag, eventVal.Field(i).String(), true
		}
	}

	return "", "", false
}
