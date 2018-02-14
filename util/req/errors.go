//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

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
