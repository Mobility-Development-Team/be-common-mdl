package strutil

import (
	"fmt"
	"reflect"
)

// StrOrEmpty Returns the string pointed by a string pointer str,
// if str is nil, return empty string "" instead
func StrOrEmpty(str *string) string {
	if str != nil {
		return *str
	}
	return ""
}

// NewPtr Copies str, and returns a pointer
func NewPtr(str string) *string {
	return &str
}

func StrOrNotProvided(str string) string {
	if str == "" {
		return "Not provided"
	}
	return str
}

func StrOrEmptyFromInterface(obj interface{}) string {
	if obj == nil {
		return ""
	}
	value := reflect.Indirect(reflect.ValueOf(obj))
	if !value.IsValid() {
		// Value is invalid, e.g. pointer of nil
		return ""
	}
	return fmt.Sprint(value.Interface())
}
