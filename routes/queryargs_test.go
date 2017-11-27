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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	QueryArgTestInt       = QueryArg{"testInt", QueryArgInt{}}
	QueryArgTestUint      = QueryArg{"testUint", QueryArgUint{}}
	QueryArgTestString    = QueryArg{"testString", QueryArgString{}}
	QueryArgTestIntSlice  = QueryArg{"testIntSlice", QueryArgIntSlice{}}
	QueryArgTestUintSlice = QueryArg{"testUintSlice", QueryArgUintSlice{}}
)

func argHandler(r *http.Request, a Arguments) (int, interface{}) {
	return 200, a
}

func intSliceIsEqual(first, second []int) bool {
	if len(first) != len(second) {
		return false
	}
	for id := range first {
		if first[id] != second[id] {
			return false
		}
	}
	return true
}

func uintSliceIsEqual(first, second []uint) bool {
	if len(first) != len(second) {
		return false
	}
	for id := range first {
		if first[id] != second[id] {
			return false
		}
	}
	return true
}

func sliceIsEqual(first, second interface{}) bool {
	switch first.(type) {
	case []int:
		return intSliceIsEqual(first.([]int), second.([]int))
	case []uint:
		return uintSliceIsEqual(first.([]uint), second.([]uint))
	default:
		return false
	}
}

func TestOverflowIntArg(t *testing.T) {
	h := ApplyDecorators(
		baseIntermediate(argHandler),
		RequireMethod{"GET"},
		WithQueryArg{QueryArgTestInt},
	)
	overflowInt64Str := "9223372036854775808"
	request := httptest.NewRequest("GET", "/test?testInt="+overflowInt64Str, nil)
	response := httptest.NewRecorder()
	status, body := h(response, request, Arguments{})
	if status != 400 {
		t.Errorf("Expected 400. Got %d (%s)", status, body)
	}
	if errorBody, ok := body.(ErrorBody); !ok {
		t.Errorf("Expected ErrorBody.")
	} else if errorBody.Error != "argument \"testInt\" must be an int" {
		t.Errorf("Expected (%v). Got (%v)", "argument \"testInt\" must be an int", errorBody.Error)
	}
}

func TestRightUintArg(t *testing.T) {
	h := ApplyDecorators(
		baseIntermediate(argHandler),
		RequireMethod{"GET"},
		WithQueryArg{QueryArgTestUint},
	)
	request := httptest.NewRequest("GET", "/test?testUint=84", nil)
	response := httptest.NewRecorder()
	status, body := h(response, request, Arguments{})
	if status != 200 {
		t.Errorf("Expected 200. Got %d (%s)", status, body)
	} else if args, ok := body.(Arguments); !ok {
		t.Errorf("Expected type Arguments.")
	} else if testUint, ok := args[QueryArgTestUint]; !ok {
		t.Errorf("testUint not in the arguments.")
	} else if testUint.(uint) != 84 {
		t.Errorf("testUint: Expected 84. Got %v", testUint)
	}
}

func TestNegativeUintArg(t *testing.T) {
	h := ApplyDecorators(
		baseIntermediate(argHandler),
		RequireMethod{"GET"},
		WithQueryArg{QueryArgTestUint},
	)
	request := httptest.NewRequest("GET", "/test?testUint=-21", nil)
	response := httptest.NewRecorder()
	status, body := h(response, request, Arguments{})
	if status != 400 {
		t.Errorf("Expected 400. Got %d (%s)", status, body)
	}
	if errorBody, ok := body.(ErrorBody); !ok {
		t.Errorf("Expected ErrorBody.")
	} else if errorBody.Error != "argument \"testUint\" must be an uint" {
		t.Errorf("Expected (%v). Got (%v)", "argument \"testUint\" must be an uint", errorBody.Error)
	}
}

func TestMultipleArg(t *testing.T) {
	h := ApplyDecorators(
		baseIntermediate(argHandler),
		RequireMethod{"GET"},
		WithQueryArg{
			QueryArgTestInt,
			QueryArgTestUint,
			QueryArgTestString,
			QueryArgTestIntSlice,
			QueryArgTestUintSlice,
		},
	)
	paramsURL := []string{
		"?testInt=-84",
		"&testUint=21",
		"&testString=test1,test2",
		"&testIntSlice=-21,-42,1",
		"&testUintSlice=21,0",
	}
	slices := []interface{}{
		int(-84),
		uint(21),
		string("test1,test2"),
		[]int{-21, -42, 1},
		[]uint{21, 0},
	}
	request := httptest.NewRequest("GET", "/test"+strings.Join(paramsURL, ""), nil)
	response := httptest.NewRecorder()
	status, body := h(response, request, Arguments{})
	if status != 200 {
		t.Errorf("Expected 200. Got %d (%s)", status, body)
	} else if args, ok := body.(Arguments); !ok {
		t.Errorf("Expected type Arguments.")
	} else {
		for id, name := range []QueryArg{
			QueryArgTestInt,
			QueryArgTestUint,
			QueryArgTestString,
			QueryArgTestIntSlice,
			QueryArgTestUintSlice,
		} {
			if slice, ok := args[name]; !ok {
				t.Errorf("%s not in the arguments.", name)
			} else {
				if id < 3 && slice != slices[id] {
					t.Errorf("Expected %v. Got %v", slice, slices[id])
				} else if id >= 3 && !sliceIsEqual(slice, slices[id]) {
					t.Errorf("Expected %v. Got %v", slice, slices[id])
				}
			}
		}
	}
}
