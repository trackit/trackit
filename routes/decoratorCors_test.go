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
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
	"testing/quick"
)

func headerStringSlice(ss []string) []string { return []string{strings.Join(ss, ",")} }
func headerBool(b bool) []string             { return []string{fmt.Sprintf("%t", b)} }

func TestCorsAllowOrigin(t *testing.T) {
	f := func(headers []string) []string {
		h := H(getFoo).With(Cors{AllowOrigin: headers})
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		s, r := h.Func(response, request, nil)
		if s != http.StatusOK {
			t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
		}
		if rt, ok := r.(string); ok {
			if rt != getFooResponse {
				t.Errorf(
					"Response should be '%s', is '%s' instead.",
					getFooResponse,
					rt,
				)
			}
		} else {
			t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", ErrUnsupportedContentType, r)
		}
		return response.Header()["Access-Control-Allow-Origin"]
	}
	if err := quick.CheckEqual(f, headerStringSlice, nil); err != nil {
		t.Error(err)
	}
}

func TestCorsAllowHeaders(t *testing.T) {
	f := func(headers []string) []string {
		h := H(getFoo).With(Cors{AllowHeaders: headers})
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		s, r := h.Func(response, request, nil)
		if s != http.StatusOK {
			t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
		}
		if rt, ok := r.(string); ok {
			if rt != getFooResponse {
				t.Errorf(
					"Response should be '%s', is '%s' instead.",
					getFooResponse,
					rt,
				)
			}
		} else {
			t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", ErrUnsupportedContentType, r)
		}
		return response.Header()["Access-Control-Allow-Headers"]
	}
	if err := quick.CheckEqual(f, headerStringSlice, nil); err != nil {
		t.Error(err)
	}
}

func TestCorsAllowCredentials(t *testing.T) {
	f := func(credentials bool) []string {
		h := H(getFoo).With(Cors{AllowCredentials: credentials})
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		s, r := h.Func(response, request, nil)
		if s != http.StatusOK {
			t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
		}
		if rt, ok := r.(string); ok {
			if rt != getFooResponse {
				t.Errorf(
					"Response should be '%s', is '%s' instead.",
					getFooResponse,
					rt,
				)
			}
		} else {
			t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", ErrUnsupportedContentType, r)
		}
		return response.Header()["Access-Control-Allow-Credentials"]
	}
	if err := quick.CheckEqual(f, headerBool, nil); err != nil {
		t.Error(err)
	}
}

func TestCorsMethods(t *testing.T) {
	h := MethodMuxer{
		http.MethodGet:  H(getFoo),
		http.MethodPost: H(postFoo),
	}.H().With(Cors{})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	s, r := h.Func(response, request, nil)
	if s != http.StatusOK {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
	}
	if rt, ok := r.(string); ok {
		if rt != getFooResponse {
			t.Errorf(
				"Response should be '%s', is '%s' instead.",
				getFooResponse,
				rt,
			)
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", ErrUnsupportedContentType, r)
	}
	expected := []string{"GET", "POST", "OPTIONS"}
	results := response.Header()["Access-Control-Allow-Methods"]
	if len(results) != 1 {
		t.Errorf("Header Access-Control-Allow-Methods should appear once. Appears %d times.", len(results))
	} else {
		result := strings.Split(results[0], ",")
		sort.Strings(expected)
		sort.Strings(result)
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Allowed methods should be %v, is %v instead.", expected, result)
		}
	}
}

func TestCorsDocumentation(t *testing.T) {
	var testCorsDocumentationExpected = HandlerDocumentation{
		Components: map[string]HandlerDocumentation{
			"method:GET":  HandlerDocumentation{},
			"method:POST": HandlerDocumentation{},
		},
	}
	h := MethodMuxer{
		http.MethodGet:  H(getFoo),
		http.MethodPost: H(postFoo),
	}.H().With(Cors{})
	request := httptest.NewRequest(http.MethodOptions, "/", nil)
	response := httptest.NewRecorder()
	s, r := h.Func(response, request, nil)
	if s != http.StatusOK {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
	}
	if rt, ok := r.(HandlerDocumentation); ok {
		if !reflect.DeepEqual(rt, testCorsDocumentationExpected) {
			t.Errorf(
				"Response should be '%#v', is '%#v' instead.",
				testCorsDocumentationExpected,
				rt,
			)
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", ErrUnsupportedContentType, r)
	}
	expected := []string{"GET", "POST", "OPTIONS"}
	results := response.Header()["Access-Control-Allow-Methods"]
	if len(results) != 1 {
		t.Errorf("Header Access-Control-Allow-Methods should appear once. Appears %d times.", len(results))
	} else {
		result := strings.Split(results[0], ",")
		sort.Strings(expected)
		sort.Strings(result)
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Allowed methods should be %v, is %v instead.", expected, result)
		}
	}
}
