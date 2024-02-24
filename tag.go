package streamline

import (
	"reflect"
)

func TagValue(inType reflect.Type) (value string, ok bool) {
	for i := 0; i < inType.NumField(); i++ {
		val := inType.Field(i).Tag.Get("streamline")
		if val != "" {
			return val, true
		}
	}

	return "", false
}

func TagFieldValue(event Event) (value string, ok bool) {
	eventVal := reflect.ValueOf(event)
	for i := 0; i < eventVal.NumField(); i++ {
		_, ok := eventVal.Type().Field(i).Tag.Lookup("streamline")
		if ok {
			// return fmt.Sprintf("%v", eventVal.Field(i).Interface()), true
			return eventVal.Field(i).String(), true
		}
	}

	return "", false
}
