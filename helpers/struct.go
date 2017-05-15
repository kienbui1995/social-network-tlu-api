package helpers

import (
	"errors"
	"reflect"
)

// IsZeroOfUnderlyingType return wether x is the is
// the zero-value of its underlying type.
func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// Replace replaces all fields of struct b that have a
// zero-value with the corresponding field value from a.
// b must be a pointer to a struct.
func Replace(a, b interface{}) error {
	// Check a.
	va := reflect.ValueOf(a)
	if va.Kind() != reflect.Struct {
		return errors.New("a is not a struct")
	}
	// Check b.
	vb := reflect.ValueOf(b)
	if vb.Kind() != reflect.Ptr {
		return errors.New("b is not a pointer")
	}
	// vb is a pointer, indirect it to get the
	// underlying value, and make sure it is a struct.
	vb = vb.Elem()
	if vb.Kind() != reflect.Struct {
		return errors.New("b is not a struct")
	}
	for i := 0; i < vb.NumField(); i++ {
		field := vb.Field(i)
		if field.CanInterface() && IsZeroOfUnderlyingType(field.Interface()) {
			// This field have a zero-value.
			// Search in a for a field with the same name.
			name := vb.Type().Field(i).Name
			fa := va.FieldByName(name)
			if fa.IsValid() {
				// Field with name was found in struct a,
				// assign its value to the field in b.
				if field.CanSet() {
					field.Set(fa)
				}
			}
		}
	}
	return nil
}
