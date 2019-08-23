package arg

import (
	"encoding"
	"reflect"

	scalar "github.com/alexflint/go-scalar"
)

var textUnmarshalerType = reflect.TypeOf([]encoding.TextUnmarshaler{}).Elem()
var argUnmarshalerType = reflect.TypeOf([]ArgUnmarshaler{}).Elem()

func canParseWrapped(t reflect.Type) bool {
	if t.Implements(argUnmarshalerType) || reflect.PtrTo(t).Implements(argUnmarshalerType) {
		return true
	}
	return scalar.CanParse(t)
}

// canParse returns true if the type can be parsed from a string
func canParse(t reflect.Type) (parseable, boolean, multiple bool) {
	parseable = canParseWrapped(t)
	boolean = isBoolean(t)
	if parseable {
		return
	}

	// Look inside pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// Look inside slice types
	if t.Kind() == reflect.Slice {
		multiple = true
		t = t.Elem()
	}

	parseable = canParseWrapped(t)
	boolean = isBoolean(t)
	if parseable {
		return
	}

	// Look inside pointer types (again, in case of []*Type)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	parseable = canParseWrapped(t)
	boolean = isBoolean(t)
	if parseable {
		return
	}

	return false, false, false
}

// isBoolean returns true if the type can be parsed from a single string
func isBoolean(t reflect.Type) bool {
	switch {
	case t.Implements(argUnmarshalerType):
		return false
	case t.Implements(textUnmarshalerType):
		return false
	case t.Kind() == reflect.Bool:
		return true
	case t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Bool:
		return true
	default:
		return false
	}
}
