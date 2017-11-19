package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testDocumentedHandlerWithRequestContentTypeExpected = `{"summary":"Get yourself some foo","description":"The route gives you some foo.","tags":{"required:oneof:contenttype":["application/json","application/csv"]}}`

func TestDocumentedHandlerDocumentationWithRequestContentType(t *testing.T) {
	h := H(getFoo).With(
		Documentation{
			Summary:     "Get yourself some foo",
			Description: "The route gives you some foo.",
		},
		RequestContentType{"application/json", "application/csv"},
	)
	bytes, err := json.Marshal(h.Documentation)
	if err == nil && string(bytes) != testDocumentedHandlerWithRequestContentTypeExpected {
		t.Errorf("JSON documentation should be '%s', is '%s' instead.", testDocumentedHandlerWithRequestContentTypeExpected, string(bytes))
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
		t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
	}
	if rt, ok := r.(error); ok {
		if rt != ErrUnsupportedContentType {
			t.Errorf("Response should be '%s', is '%s' instead.", ErrUnsupportedContentType.Error(), rt.Error())
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", ErrUnsupportedContentType, r)
	}
}
