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
	"testing"
)

type helloStruct struct {
	Hello string `json:"hello"`
}

func helloHandler(r *http.Request, a Arguments) (int, interface{}) {
	return 200, helloStruct{
		Hello: "world",
	}
}

func TestRequireMethodAcceptsAllowed(t *testing.T) {
	d := RequireMethod{"GET", "POST"}
	h := ApplyDecorators(
		baseIntermediate(helloHandler),
		d,
	)
	for _, m := range d {
		request := httptest.NewRequest(m, "/hello", nil)
		response := httptest.NewRecorder()
		status, body := h(response, request, nil)
		if status != 200 {
			t.Errorf("Expected 200 on %s. Got %d.", m, status)
		}
		if _, ok := body.(helloStruct); ok == false {
			t.Errorf("Expected helloStruct as response, got %T %#v.", body, body)
		}
	}
}

func TestRequireMethodNotAcceptedRejected(t *testing.T) {
	d := RequireMethod{"GET", "POST"}
	h := ApplyDecorators(
		baseIntermediate(helloHandler),
		d,
	)
	for _, m := range [...]string{"PUT", "PATCH", "DELETE"} {
		request := httptest.NewRequest(m, "/hello", nil)
		response := httptest.NewRecorder()
		status, body := h(response, request, nil)
		if status != 405 {
			t.Errorf("Expected 405 on %s. Got %d.", m, status)
		}
		if body != methodNotAllowed {
			t.Errorf("Expected methodNotAllowed as response, got %T %#v.", body, body)
		}
	}
}

func TestRequireMethodHasAllowedHeader0(t *testing.T) {
	d := RequireMethod{"GET", "POST"}
	h := ApplyDecorators(
		baseIntermediate(helloHandler),
		d,
	)
	request := httptest.NewRequest("OPTIONS", "/hello", nil)
	response := httptest.NewRecorder()
	status, _ := h(response, request, nil)
	if status != 200 {
		t.Errorf("Expected 200 on OPTIONS. Got %d.", status)
	}
	if a := response.Header()["Allow"][0]; a != "GET,POST,OPTIONS" {
		t.Errorf("Expected 'GET,POST,OPTIONS' Allow header. Got '%s'.", a)
	}
}

func TestRequireMethodHasAllowedHeader1(t *testing.T) {
	d := RequireMethod{}
	h := ApplyDecorators(
		baseIntermediate(helloHandler),
		d,
	)
	request := httptest.NewRequest("OPTIONS", "/hello", nil)
	response := httptest.NewRecorder()
	status, _ := h(response, request, nil)
	if status != 200 {
		t.Errorf("Expected 200 on OPTIONS. Got %d.", status)
	}
	if a := response.Header()["Allow"][0]; a != "OPTIONS" {
		t.Errorf("Expected 'OPTIONS' Allow header. Got '%s'.", a)
	}
}
