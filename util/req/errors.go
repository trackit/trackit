package req

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidationError is an error representing one or more validation errors.
type ValidationError []error

func (ve ValidationError) Error() string {
	es := make([]string, len(ve))
	for i := range es {
		es[i] = ve[i].Error()
	}
	return strings.Join(es, "; ")
}

// FieldValidationError signals that a field failed to validate.
type FieldValidationError struct {
	Field reflect.StructField
	Err   error
}

func (fve FieldValidationError) Error() string {
	n := getJsonName(fve.Field)
	return fmt.Sprintf("%s: %s", n, fve.Err.Error())
}

// getJsonName retrieves the name package encoding/json would associate to a
// struct field.
func getJsonName(fld reflect.StructField) string {
	jsonTag := fld.Tag.Get("json")
	jsonTagParts := strings.SplitN(jsonTag, ",", 2)
	jsonName := jsonTagParts[0]
	if jsonName == "" {
		return fld.Name
	} else {
		return jsonName
	}
}
