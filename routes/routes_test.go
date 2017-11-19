package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const getFooResponse = "HAVE SOME FOO"

func getFoo(r *http.Request, a Arguments) (int, interface{}) {
	return http.StatusOK, getFooResponse
}

const postFooResponse = "GIMME YOUR FOO"

func postFoo(r *http.Request, a Arguments) (int, interface{}) {
	return http.StatusOK, postFooResponse
}

const testDocumentedHandlerExpected = `{"summary":"Get yourself some foo","description":"The route gives you some foo.","tags":{"foo":["yes"]}}`

func TestDocumentedHandlerDocumentation(t *testing.T) {
	h := H(getFoo).With(
		Documentation{
			Summary:     "Get yourself some foo",
			Description: "The route gives you some foo.",
			Tags:        Tags{"foo": []string{"yes"}},
		},
	)
	bytes, err := json.Marshal(h.Documentation)
	if err == nil && string(bytes) != testDocumentedHandlerExpected {
		t.Errorf("JSON documentation should be '%s', is '%s' instead.", testDocumentedHandlerExpected, string(bytes))
	} else if err != nil {
		t.Errorf("Error should be nil, is '%s' instead.", err.Error())
	}
}

func TestDocumentedHandlerFunctionality(t *testing.T) {
	h := H(getFoo).With(
		Documentation{
			Summary:     "Get yourself some foo",
			Description: "The route gives you some foo.",
			Tags:        Tags{"foo": []string{"yes"}},
		},
	)
	s, r := h.Func(nil, nil, nil)
	if s != http.StatusOK {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
	}
	if rt, ok := r.(string); ok {
		if rt != getFooResponse {
			t.Errorf("Response should be '%s', is '%s' instead.", getFooResponse, rt)
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", getFooResponse, rt)
	}
}

const testRegistrationPattern = "/foo"

func TestRegistration(t *testing.T) {
	var handlerRun bool
	f := func(_ *http.Request, _ Arguments) (int, interface{}) { handlerRun = true; return 200, nil }
	h := H(f)
	h.Register(testRegistrationPattern)
	if l := len(RegisteredHandlers); l != 1 {
		t.Errorf("Count of registered handlers should be 1, is %d.", l)
	} else {
		rh := RegisteredHandlers[0]
		if rh.Pattern != testRegistrationPattern {
			t.Errorf("Pattern should be %s, is %s instead.", testRegistrationPattern, rh.Pattern)
		}
		rh.Handler.Func(nil, nil, nil)
		if !handlerRun {
			t.Error("Handler should have run, hasn't.")
		}
	}
	resetRegisteredHandlers()
}

const testMethodMuxerDocumentationExpected = `{"summary":"Interacts with foo.","components":{"GET":{"summary":"Gets foo.","componentType":"method"},"POST":{"summary":"Posts foo.","componentType":"method"}}}`

func TestMethodMuxerDocumentation(t *testing.T) {
	h := MethodMuxer{
		http.MethodGet:  H(getFoo).With(Documentation{Summary: "Gets foo."}),
		http.MethodPost: H(postFoo).With(Documentation{Summary: "Posts foo."}),
	}.H().With(Documentation{Summary: "Interacts with foo."})
	bytes, err := json.Marshal(h.Documentation)
	if err == nil && string(bytes) != testMethodMuxerDocumentationExpected {
		t.Errorf("JSON documentation should be '%s', is '%s' instead.", testMethodMuxerDocumentationExpected, string(bytes))
	} else if err != nil {
		t.Errorf("Error should be nil, is '%s' instead.", err.Error())
	}
}

func TestMethodMuxerFunctionality(t *testing.T) {
	h := MethodMuxer{
		http.MethodGet:  H(getFoo).With(Documentation{Summary: "Gets foo."}),
		http.MethodPost: H(postFoo).With(Documentation{Summary: "Posts foo."}),
	}.H().With(Documentation{Summary: "Interacts with foo."})
	s, r := h.Func(nil, httptest.NewRequest(http.MethodGet, "/", nil), nil)
	if s != http.StatusOK {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
	}
	if rt, ok := r.(string); ok {
		if rt != getFooResponse {
			t.Errorf("Response should be '%s', is '%s' instead.", getFooResponse, rt)
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", getFooResponse, rt)
	}
}

func TestMethodMuxerMethodNotAllowed(t *testing.T) {
	h := MethodMuxer{
		http.MethodGet:  H(getFoo).With(Documentation{Summary: "Gets foo."}),
		http.MethodPost: H(postFoo).With(Documentation{Summary: "Posts foo."}),
	}.H().With(Documentation{Summary: "Interacts with foo."})
	s, r := h.Func(nil, httptest.NewRequest(http.MethodPut, "/", nil), nil)
	if s != http.StatusMethodNotAllowed {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
	}
	if rt, ok := r.(error); ok {
		if rt != ErrMethodNotAllowed {
			t.Errorf("Response should be '%s', is '%s' instead.", ErrMethodNotAllowed, rt)
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", ErrMethodNotAllowed, r)
	}
}
