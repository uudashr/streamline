package streamline

import (
	"reflect"
)

const (
	tagName = "streamline"
)

func TagValue(inType reflect.Type) (value string, ok bool) {
	for i := 0; i < inType.NumField(); i++ {
		val := inType.Field(i).Tag.Get(tagName)
		if val != "" {
			return val, true
		}
	}

	return "", false
}

func TagFieldValue(event Event) (value string, ok bool) {
	_, val, ok := Tag(event)
	return val, ok
}

func Tag(event Event) (tag string, value string, ok bool) {
	eventVal := reflect.ValueOf(event)
	for i := 0; i < eventVal.NumField(); i++ {
		tag, ok := eventVal.Type().Field(i).Tag.Lookup(tagName)
		if ok {
			return tag, eventVal.Field(i).String(), true
		}
	}

	return "", "", false
}
