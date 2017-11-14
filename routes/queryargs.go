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

package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type (
	// QueryArgInt is a type allowing to know the expected int
	// type for the argument. QueryArgInt can be interfaced to
	// QueryArgType.
	QueryArgInt struct{}

	// QueryArgUint is a type allowing to know the expected uint
	// type for the argument. QueryArgUint can be interfaced to
	// QueryArgType.
	QueryArgUint struct{}

	// QueryArgString is a type allowing to know the expected string
	// type for the argument. QueryArgString can be interfaced to
	// QueryArgType.
	QueryArgString struct{}

	// QueryArgIntSlice is a type allowing to know the expected []int
	// type for the argument. QueryArgIntSlice can be interfaced to
	// QueryArgType.
	QueryArgIntSlice struct{}

	// QueryArgUintSlice is a type allowing to know the expected []uint
	// type for the argument. QueryArgUintSlice can be interfaced to
	// QueryArgType.
	QueryArgUintSlice struct{}

	// QueryArgStringSlice is a type allowing to know the expected []string
	// type for the argument. QueryArgStringSlice can be interfaced to
	// QueryArgType.
	QueryArgStringSlice struct{}

	// QueryArgType is an interface used by all the type above. The Parse
	// method takes the raw value of the argument and returns its typed value.
	// An error can be returned if the value could not be parse. The error's
	// message contains %s which has to be replaced by the argument's name
	// before being displayed.
	QueryArgType interface {
		Parse(string) (interface{}, error)
	}

	// QueryArg is a structure defining an argument by its name and its type.
	QueryArg struct {
		Name string
		Type QueryArgType
	}

	// WithQueryArg is a slice of QueryArg. That is the type that contains all
	// the arguments to parse in the URL. WithQueryArg has a method Decorate
	// called to apply the decorators on an endpoint.
	WithQueryArg []QueryArg
)

// Parse is the method of QueryArgInt allowing to type the row value to
// an int value. It can return an error, the error's message contains %s
// which has to be replaced by the argument's name before being displayed.
func (d QueryArgInt) Parse(val string) (interface{}, error) {
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return i, nil
	}
	return nil, errors.New("argument \"%s\" must be an int")
}

// Parse is the method of QueryArgUint allowing to type the row value to
// an uint value. It can return an error, the error's message contains %s
// which has to be replaced by the argument's name before being displayed.
func (d QueryArgUint) Parse(val string) (interface{}, error) {
	if i, err := strconv.ParseUint(val, 10, 64); err == nil {
		return i, nil
	}
	return nil, errors.New("argument \"%s\" must be an uint")
}

// Parse is the method of QueryArgString allowing to type the row value to
// a string value. It can't return an error since the row value is already
// string typed.
func (d QueryArgString) Parse(val string) (interface{}, error) {
	return val, nil
}

// Parse is the method of QueryArgIntSlice allowing to type the row value to
// a []int value. It can return an error, the error's message contains %s
// which has to be replaced by the argument's name before being displayed.
func (d QueryArgIntSlice) Parse(val string) (interface{}, error) {
	vals := strings.Split(val, ",")
	res := make([]int64, 0, len(vals))
	for _, v := range vals {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			res = append(res, i)
		} else {
			return nil, errors.New("argument \"%s\" must be a slice of int")
		}
	}
	return res, nil
}

// Parse is the method of QueryArgUintSlice allowing to type the row value to
// a []uint value. It can return an error, the error's message contains %s
// which has to be replaced by the argument's name before being displayed.
func (d QueryArgUintSlice) Parse(val string) (interface{}, error) {
	vals := strings.Split(val, ",")
	res := make([]uint64, 0, len(vals))
	for _, v := range vals {
		if i, err := strconv.ParseUint(v, 10, 64); err == nil {
			res = append(res, i)
		} else {
			return nil, errors.New("argument \"%s\" must be a slice of uint")
		}
	}
	return res, nil
}

// Parse is the method of QueryArgStringSlice allowing to type the row value to
// a []string value. It can return an error, the error's message contains %s
// which has to be replaced by the argument's name before being displayed.
func (d QueryArgStringSlice) Parse(val string) (interface{}, error) {
	res := strings.Split(val, ",")
	if len(res) > 0 {
		return res, nil
	}
	return nil, errors.New("argument \"%s\" must be a slice of string")
}

// Decorate is the function called to apply the decorators to an endpoint. It returns
// a function. This function produces a 400 error code with a json error message or
// calls the next IntermediateHandler.
// The goal of this function is to get the URL parameters to store them in
// the Arguments.
//
// To use this decorator, create a var:
// 		QueryArgProductsName = QueryArg{"productsName", QueryArgStringSlice{}}
//
// In this example, we need a slice of string called productsName that will be
// our products name so we indicate to the decorator with QueryArgStringSlice{}.
//
// Then, register your endpoint as usual:
//			Register(
//				"/products/cost",
//				getProductsCost,
//				RequireMethod{"GET"},
//				RequireContentType{"application/json"},
//				db.WithTransaction{db.Db},
//				users.WithAuthenticatedUser{},
//				WithQueryArg{QueryArgProductsName},
//			)
//
// See the last decorator WithQueryArg{QueryArgProductsName}, it will store
// a slice of product name in the arguments.
//
// Thus, if the user calls
//		/products/cost?productsName=EC2,ES,RDS
// you will be able to get the productsName with the key productsName in the
// arguments. The value will be ["EC2", "ES", "RDS"].
func (d WithQueryArg) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (status int, output interface{}) {
		for _, arg := range d {
			if rawVal := r.URL.Query().Get(arg.Name); rawVal != "" {
				if val, err := arg.Type.Parse(rawVal); err == nil {
					a[arg.Name] = val
				} else {
					msg := fmt.Sprintf(err.Error(), arg.Name)
					return 400, ErrorBody{msg}
				}
			} else {
				msg := fmt.Sprintf("argument \"%s\" not found", arg.Name)
				return 400, ErrorBody{msg}
			}
		}
		return h(w, r, a)
	}
}
