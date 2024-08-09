package funcs

import (
	"reflect"
)

// IsAvailable return a boolean value indicating whether 'name' is a named property of 'data'.
func IsAvailable(name string, data interface{}) bool {

	// https://stackoverflow.com/questions/44675087/golang-template-variable-isset

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	return v.FieldByName(name).IsValid()
}
