package jsonutil

import (
	"encoding/json"
	"fmt"
	"reflect"

	logger "github.com/sirupsen/logrus"
)

// Convert a json compatible struct/map into another object.
//
// It is interally done by marshalling src into json first,
// then unmarshal it into another object
func MarshalInto(src, dest interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("unable to marshal object type %v as json: %w", reflect.TypeOf(src), err)
	}
	if err := json.Unmarshal(b, dest); err != nil {
		return fmt.Errorf("unable to marshal object type %v into %v: %w", reflect.TypeOf(src), reflect.TypeOf(dest), err)
	}
	return nil
}

// Same as MarshalInto but logs the error as warning instead of returning it
func ShouldMarshalInto(src, dest interface{}) {
	if err := MarshalInto(src, dest); err != nil {
		logger.Warn("[ShouldMarshalInto] Unable to marshal types: ", err)
	}
}
