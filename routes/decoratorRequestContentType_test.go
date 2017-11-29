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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {

}

const testDocumentedHandlerWithRequestContentTypeExpected = `{
	"summary": "Get yourself some foo",
	"description": "The route gives you some foo.",
	"tags": {
		"require:contenttype": [
			"application/json",
			"application/csv"
		]
	}
}`

func TestDocumentedHandlerDocumentationWithRequestContentType(t *testing.T) {
	h := H(getFoo).With(
		Documentation{
			Summary:     "Get yourself some foo",
			Description: "The route gives you some foo.",
		},
		RequestContentType{"application/json", "application/csv"},
	)
	bytes, err := json.MarshalIndent(h.Documentation, "", "\t")
	if err == nil && string(bytes) != testDocumentedHandlerWithRequestContentTypeExpected {
		t.Errorf(
			"JSON documentation should be '%s', is '%s' instead.",
			testDocumentedHandlerWithRequestContentTypeExpected,
			string(bytes),
		)
	} else if err != nil {
		t.Errorf("Error should be nil, is '%s' instead.", err.Error())
	}
}

func TestDocumentedHandlerWithGoodRequestContentType(t *testing.T) {
	h := H(getFoo).With(
		Documentation{
			Summary:     "Get yourself some foo",
			Description: "The route gives you some foo.",
		},
		RequestContentType{"application/json", "application/csv"},
	)
	request := httptest.NewRequest(http.MethodPut, "/", nil)
	request.Header["Content-Type"] = []string{"application/json"}
	s, r := h.Func(nil, request, nil)
	if s != http.StatusOK {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
	}
	if rt, ok := r.(string); ok {
		if rt != getFooResponse {
			t.Errorf("Response should be '%s', is '%s' instead.", getFooResponse, rt)
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", getFooResponse, r)
	}
}

func TestDocumentedHandlerWithBadRequestContentType(t *testing.T) {
	h := H(getFoo).With(
		Documentation{
			Summary:     "Get yourself some foo",
			Description: "The route gives you some foo.",
		},
		RequestContentType{"application/json", "application/csv"},
	)
	request := httptest.NewRequest(http.MethodPut, "/", nil)
	request.Header["Content-Type"] = []string{"text/html"}
	s, r := h.Func(nil, request, nil)
	if s != http.StatusBadRequest {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusBadRequest, s)
	}
	if rt, ok := r.(error); ok {
		if rt != ErrUnsupportedContentType {
			t.Errorf(
				"Response should be '%s', is '%s' instead.",
				ErrUnsupportedContentType.Error(),
				rt.Error(),
			)
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", ErrUnsupportedContentType, r)
	}
}

func TestIgnoredMethod(t *testing.T) {
	for _, method := range [...]string{http.MethodGet, http.MethodHead, http.MethodDelete, http.MethodOptions} {
		h := H(getFoo).With(
			Documentation{
				Summary:     "Get yourself some foo",
				Description: "The route gives you some foo.",
			},
			RequestContentType{"application/json", "application/csv"},
		)
		request := httptest.NewRequest(method, "/", nil)
		s, r := h.Func(nil, request, nil)
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
			t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", getFooResponse, r)
		}
	}
}

func TestNoRequestContentType(t *testing.T) {
	h := H(getFoo).With(
		Documentation{
			Summary:     "Get yourself some foo",
			Description: "The route gives you some foo.",
		},
		RequestContentType{"application/json", "application/csv"},
	)
	request := httptest.NewRequest(http.MethodPut, "/", nil)
	request.Header["Content-Type"] = []string{}
	s, r := h.Func(nil, request, nil)
	if s != http.StatusBadRequest {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusBadRequest, s)
	}
	if rt, ok := r.(error); ok {
		if rt != ErrMissingContentType {
			t.Errorf(
				"Response should be '%s', is '%s' instead.",
				ErrMissingContentType.Error(),
				rt.Error(),
			)
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", ErrUnsupportedContentType, r)
	}
}

func TestMultipleRequestContentTypes(t *testing.T) {
	h := H(getFoo).With(
		Documentation{
			Summary:     "Get yourself some foo",
			Description: "The route gives you some foo.",
		},
		RequestContentType{"application/json", "application/csv"},
	)
	request := httptest.NewRequest(http.MethodPut, "/", nil)
	request.Header["Content-Type"] = []string{"application/json", "text/html"}
	s, r := h.Func(nil, request, nil)
	if s != http.StatusBadRequest {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusBadRequest, s)
	}
	if rt, ok := r.(error); ok {
		if rt != ErrMultipleContentTypes {
			t.Errorf(
				"Response should be '%s', is '%s' instead.",
				ErrMultipleContentTypes.Error(),
				rt.Error(),
			)
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", ErrUnsupportedContentType, r)
	}
}
