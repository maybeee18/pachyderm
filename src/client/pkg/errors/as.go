package errors

import (
	"reflect"

	"github.com/pkg/errors"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func tryAs(err error, targetVal reflect.Value) bool {
	e := targetVal.Type().Elem()
	if e.Kind() == reflect.Interface || e.Implements(errorType) {
		return errors.As(err, targetVal.Interface())
	}
	return false
}

// As finds the first error in err's chain that matches the target's type, and
// if so, sets target to that error value and returns true.
// As is a wrapper for the underlying errors.As function, which may panic or
// return unexpected results based on how err was constructed (with or without
// a pointer).  This works by inspecting the type of target and attempting
// multiple errors.As calls if necessary.
func As(err error, target interface{}) bool {
	// Check the type of target, it must be a pointer to an error, or a pointer to a pointer to an error
	v := reflect.ValueOf(target)

	switch v.Kind() {
	case reflect.Ptr:
		vp := reflect.New(v.Type())

		// Attempt unwrapping a nested pointer
		if v.Elem().Kind() == reflect.Ptr && tryAs(err, v.Elem()) {
			return true
		}

		// Attempt wrapping with an additional pointer
		if tryAs(err, vp) {
			v.Elem().Set(vp.Elem().Elem())
			return true
		}

		// Attempt the passed target as-is
		return tryAs(err, v)
	}
	return false
}
