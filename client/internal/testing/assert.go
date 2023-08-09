package testing

import (
	"reflect"
	"testing"
)

func AssertDeepEqual(t *testing.T, expected any, actual any, message string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s\nexpected: %+v\ngot: %+v", message, expected, actual)
	}
}

func AssertNil(t *testing.T, value any, message string) {
	if value != nil {
		t.Fatalf("%s\nexpected value to be nil, but was=%+v", message, value)
	}
}

func AssertNotNil(t *testing.T, value any, message string) {
	if value == nil {
		t.Fatalf("%s\nexpected value to not be nil, but it was", message)
	}
}
