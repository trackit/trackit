package req

import (
	"errors"
	"reflect"
	"strings"
)

// internalValidator validates a value, returning nil in case of success.
type internalValidator func(reflect.Value) error

// Validator validates a value, returning nil in case of success.
type Validator func(interface{}) error

const (
	// StructTagName is the name of the tags package req bases its behavior
	// on.
	StructTagName = "req"
	// StructTagNonZero indicates that a field with a zero value shall fail
	// validation.
	StructTagNonZero = "nonzero"
)

var (
	ErrUnsupportedType    = errors.New("unsupported type")
	ErrUnknownDirective   = errors.New("unknown directive")
	ErrTestNonZeroIsZero  = errors.New("value is zero")
	ErrCannotValidateType = errors.New("cannot validate type")
)

// CreateValidator creates a validator for the type of its argument.
func CreateValidator(example interface{}) (Validator, error) {
	typ := reflect.TypeOf(example)
	switch typ.Kind() {
	case reflect.Struct:
		iv, err := createValidatorForStruct(typ)
		if err == nil {
			return exposedValidator(iv, typ), nil
		} else {
			return nil, err
		}
	default:
		return nil, ErrUnsupportedType
	}
}

// exposedValidator wraps an internalValidator into a Validator suitable to
// being exported.
func exposedValidator(iv internalValidator, typ reflect.Type) Validator {
	if iv == nil {
		return nil
	} else {
		return func(in interface{}) error {
			ityp := reflect.TypeOf(in)
			if ityp != typ {
				return ErrCannotValidateType
			} else {
				return iv(reflect.ValueOf(in))
			}
		}
	}
}

// createValidatorForStruct builds a validator function for a given structure
// type.
func createValidatorForStruct(typ reflect.Type) (internalValidator, error) {
	var vs []internalValidator
	for i, fc := 0, typ.NumField(); i < fc; i++ {
		f := typ.FieldByIndex([]int{i})
		if v, err := createValidatorForField(typ, f); err != nil {
			return nil, err
		} else if v != nil {
			vs = append(vs, v)
		}
	}
	return aggregateValidators(vs), nil
}

// createValidatorForField creates a validator function for a given structure
// field.
func createValidatorForField(typ reflect.Type, fld reflect.StructField) (internalValidator, error) {
	if tag := fld.Tag.Get(StructTagName); tag == "" {
		return nil, nil
	} else {
		var tags = strings.Split(tag, ",")
		var tests []internalValidator
		for _, tag := range tags {
			if test, err := getTestFor(typ, fld, tag); err != nil {
				return nil, err
			} else if test != nil {
				tests = append(tests, test)
			}
		}
		av := aggregateValidators(tests)
		return func(val reflect.Value) error {
			val = val.FieldByIndex(fld.Index)
			return av(val)
		}, nil
	}
}

// aggregateValidators creates a single validator from a collection of
// validators. When called it will allways run all child validators, returning
// all errors.
func aggregateValidators(fvs []internalValidator) internalValidator {
	if fvs == nil || len(fvs) == 0 {
		return nil
	} else {
		return func(v reflect.Value) error {
			var errs ValidationError
			for _, fv := range fvs {
				if err := fv(v); err != nil {
					errs = append(errs, err)
				}
			}
			if errs == nil { // Must return explicit nil to actually return nil
				return nil
			} else {
				return errs
			}
		}
	}
}

// getTestFor builds a validator for a given structure field and tag.
func getTestFor(typ reflect.Type, fld reflect.StructField, tag string) (internalValidator, error) {
	switch tag {
	case StructTagNonZero:
		return testNonZero(fld), nil
	default:
		return nil, ErrUnknownDirective
	}
}

// testNonZero builds a field validator which ensures the field's value is not
// the zero value for its type.
func testNonZero(fld reflect.StructField) internalValidator {
	zero := reflect.Zero(fld.Type).Interface()
	return func(val reflect.Value) error {
		if reflect.DeepEqual(val.Interface(), zero) {
			return FieldValidationError{
				Field: fld,
				Err:   ErrTestNonZeroIsZero,
			}
		} else {
			return nil
		}
	}
}
