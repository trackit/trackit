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
	"time"
)

var (
	QueryArgTestBool           = QueryArg{"testBool", "Test boolean", QueryArgBool{}, false}
	QueryArgTestInt            = QueryArg{"testInt", "Test signed integer", QueryArgInt{}, false}
	QueryArgTestUint           = QueryArg{"testUint", "Test unsigned integer", QueryArgUint{}, false}
	QueryArgTestString         = QueryArg{"testString", "Test string", QueryArgString{}, false}
	QueryArgTestOptionalString = QueryArg{"testString", "Test string", QueryArgString{}, true}
	QueryArgTestIntSlice       = QueryArg{"testIntSlice", "Test signed integer slice", QueryArgIntSlice{}, false}
	QueryArgTestUintSlice      = QueryArg{"testUintSlice", "Test unsigned integer slice", QueryArgUintSlice{}, false}
	QueryArgTestStringSlice    = QueryArg{"testStringSlice", "Test string slice", QueryArgStringSlice{}, false}
	QueryArgTestDate           = QueryArg{"testDate", "Test date", QueryArgDate{}, false}
)

func argHandler(r *http.Request, a Arguments) (int, interface{}) {
	return http.StatusOK, a
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

func stringSliceIsEqual(first, second []string) bool {
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
	case []string:
		return stringSliceIsEqual(first.([]string), second.([]string))
	default:
		return false
	}
}

const testOverflowIntArgExpectedError = `query arg 'testInt': must be an int`

func TestGoodBool(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{
			QueryArgTestBool,
		},
	)
	paramsURL := "false"
	request := httptest.NewRequest("GET", "/test?testBool="+paramsURL, nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusOK {
		t.Errorf("Expected %d but got %d (%v)", http.StatusOK, status, body)
	} else if args, ok := body.(Arguments); !ok {
		t.Errorf("Expected type Arguments")
	} else if args[QueryArgTestBool] == true {
		t.Errorf("Expected false but got %v", args[QueryArgTestBool])
	}
}

func TestEmptyBool(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{
			QueryArgTestBool,
		},
	)
	request := httptest.NewRequest("GET", "/test?testBool", nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusOK {
		t.Errorf("Expected %d but got %d (%v)", http.StatusOK, status, body)
	} else if args, ok := body.(Arguments); !ok {
		t.Errorf("Expected type Arguments")
	} else if args[QueryArgTestBool] == false {
		t.Errorf("Expected true but got %v", args[QueryArgTestBool])
	}
}

func TestOverflowIntArg(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{QueryArgTestInt},
	)
	overflowInt64Str := "9223372036854775808"
	request := httptest.NewRequest("GET", "/test?testInt="+overflowInt64Str, nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusBadRequest {
		t.Errorf("Expected http.StatusBadRequest (400). Got %d (%s)", status, body)
	}
	if err, ok := body.(error); !ok {
		t.Errorf("Expected error.")
	} else if err.Error() != testOverflowIntArgExpectedError {
		t.Errorf("Expected (%v). Got (%v)", testOverflowIntArgExpectedError, err.Error())
	}
}

func TestRightUintArg(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{QueryArgTestUint},
	)
	request := httptest.NewRequest("GET", "/test?testUint=84", nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusOK {
		t.Errorf("Expected http.StatusOK (200). Got %d (%s)", status, body)
	} else if args, ok := body.(Arguments); !ok {
		t.Errorf("Expected type Arguments.")
	} else if testUint, ok := args[QueryArgTestUint]; !ok {
		t.Errorf("testUint not in the arguments.")
	} else if testUint.(uint) != 84 {
		t.Errorf("testUint: Expected 84. Got %v", testUint)
	}
}

func TestOptionalStringArg(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{QueryArgTestOptionalString},
	)
	request := httptest.NewRequest("GET", "/test?testString=Hi!", nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusOK {
		t.Errorf("Expected http.StatusOK (200). Got %d (%s)", status, body)
	} else if args, ok := body.(Arguments); !ok {
		t.Errorf("Expected type Arguments.")
	} else if testString, ok := args[QueryArgTestOptionalString]; !ok {
		t.Errorf("testString not in the arguments.")
	} else if testString.(string) != "Hi!" {
		t.Errorf("testString: Expected Hi!. Got %v", testString)
	}

	request = httptest.NewRequest("GET", "/test", nil)
	response = httptest.NewRecorder()
	status, body = h.Func(response, request, Arguments{})
	if status != http.StatusOK {
		t.Errorf("Expected http.StatusOK (200). Got %d (%s)", status, body)
	} else if args, ok := body.(Arguments); !ok {
		t.Errorf("Expected type Arguments.")
	} else if _, ok := args[QueryArgTestOptionalString]; ok {
		t.Errorf("testString in the arguments.")
	}
}

const testNegativeUintArgExpectedError = `query arg 'testUint': must be a uint`

func TestNegativeUintArg(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{QueryArgTestUint},
	)
	request := httptest.NewRequest("GET", "/test?testUint=-21", nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusBadRequest {
		t.Errorf("Expected http.StatusBadRequest (400). Got %d (%s)", status, body)
	}
	if err, ok := body.(error); !ok {
		t.Errorf("Expected error.")
	} else if err.Error() != testNegativeUintArgExpectedError {
		t.Errorf("Expected (%v). Got (%v)", testNegativeUintArgExpectedError, err.Error())
	}
}

func TestMultipleArg(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{
			QueryArgTestInt,
			QueryArgTestUint,
			QueryArgTestString,
			QueryArgTestIntSlice,
			QueryArgTestUintSlice,
		},
	)
	paramsURL := []string{
		"testInt=-84",
		"testUint=21",
		"testString=test1,test2",
		"testIntSlice=-21,-42,1",
		"testUintSlice=21,0",
	}
	slices := []interface{}{
		int(-84),
		uint(21),
		string("test1,test2"),
		[]int{-21, -42, 1},
		[]uint{21, 0},
	}
	request := httptest.NewRequest("GET", "/test?"+strings.Join(paramsURL, "&"), nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusOK {
		t.Errorf("Expected http.StatusOK (200). Got %d (%v)", status, body)
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
				t.Errorf("%s not in the arguments.", name.Name)
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

func TestMissingIntSlice(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{
			QueryArgTestIntSlice,
		},
	)
	paramsURL := []string{
		"testString=test",
	}
	request := httptest.NewRequest("GET", "/test?"+strings.Join(paramsURL, "&"), nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusBadRequest {
		t.Errorf("Expected %d. Got %d (%v)", http.StatusBadRequest, status, body)
	} else if _, ok := body.(error); !ok {
		t.Errorf("Expected error.")
	}
}

func TestBadIntSlice(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{
			QueryArgTestIntSlice,
		},
	)
	paramsURL := []string{
		"testIntSlice=test",
	}
	request := httptest.NewRequest("GET", "/test?"+strings.Join(paramsURL, "&"), nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusBadRequest {
		t.Errorf("Expected %d. Got %d (%v)", http.StatusBadRequest, status, body)
	} else if _, ok := body.(error); !ok {
		t.Errorf("Expected error.")
	}
}

func TestBadUintSlice(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{
			QueryArgTestUintSlice,
		},
	)
	paramsURL := []string{
		"testUintSlice=test",
	}
	request := httptest.NewRequest("GET", "/test?"+strings.Join(paramsURL, "&"), nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusBadRequest {
		t.Errorf("Expected %d. Got %d (%v)", http.StatusBadRequest, status, body)
	} else if _, ok := body.(error); !ok {
		t.Errorf("Expected error.")
	}
}

func TestGoodStringSlice(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{
			QueryArgTestStringSlice,
		},
	)
	paramsURL := []string{
		"testStringSlice=foo,bar",
	}
	expectedResult := []string{
		"foo",
		"bar",
	}
	request := httptest.NewRequest("GET", "/test?"+strings.Join(paramsURL, "&"), nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusOK {
		t.Errorf("Expected %d but got %d (%v)", http.StatusOK, status, body)
	} else if args, ok := body.(Arguments); !ok {
		t.Errorf("Expected type Arguments")
	} else if !sliceIsEqual(expectedResult, args[QueryArgTestStringSlice]) {
		t.Errorf("Expected %v but got %v", expectedResult, args[QueryArgTestStringSlice])
	}
}

func TestEmptyStringSlice(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{
			QueryArgTestStringSlice,
		},
	)
	paramsURL := []string{
		"testStringSlice=",
	}
	request := httptest.NewRequest("GET", "/test?"+strings.Join(paramsURL, "&"), nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusBadRequest {
		t.Errorf("Expected %d but got %d (%v)", http.StatusBadRequest, status, body)
	} else if _, ok := body.(error); !ok {
		t.Errorf("Expected error.")
	}
}

func TestGoodDate(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{
			QueryArgTestDate,
		},
	)
	paramsURL := []string{
		"testDate=2017-12-21",
	}
	expectedDate, _ := time.Parse(iso8601DateFormat, "2017-12-21")
	request := httptest.NewRequest("GET", "/test?"+strings.Join(paramsURL, "&"), nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusOK {
		t.Errorf("Expected %d but got %d (%v)", http.StatusOK, status, body)
	} else if args, ok := body.(Arguments); !ok {
		t.Errorf("Expected type Arguments")
	} else if responseTime, ok := args[QueryArgTestDate].(time.Time); !ok {
		t.Errorf("Expected type time")
	} else if !expectedDate.Equal(responseTime) {
		t.Errorf("Expected time %v but got time %v", expectedDate, responseTime)
	}
}

func TestBadDate(t *testing.T) {
	h := H(argHandler).With(
		QueryArgs{
			QueryArgTestDate,
		},
	)
	paramsURL := []string{
		"testDate=2017-31-42:foo",
	}
	request := httptest.NewRequest("GET", "/test?"+strings.Join(paramsURL, "&"), nil)
	response := httptest.NewRecorder()
	status, body := h.Func(response, request, Arguments{})
	if status != http.StatusBadRequest {
		t.Errorf("Expected %d but got %d (%v)", http.StatusBadRequest, status, body)
	} else if _, ok := body.(error); !ok {
		t.Errorf("Expected error.")
	}
}
