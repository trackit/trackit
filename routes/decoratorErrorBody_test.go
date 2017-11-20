package routes

import (
	"errors"
	"net/http"
	"testing"
)

const testErrorText = "test error"

func erroringHandler(r *http.Request, a Arguments) (int, interface{}) {
	return http.StatusInternalServerError, errors.New(testErrorText)
}

func TestErrorBodyWrapsErrors(t *testing.T) {
	h := H(erroringHandler).With(ErrorBody{})
	s, r := h.Func(nil, nil, nil)
	if s != http.StatusInternalServerError {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
	}
	if rt, ok := r.(errorBody); ok {
		if rt.Error != testErrorText {
			t.Errorf("Response should be '%s', is '%s' instead.", testErrorText, rt)
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", errorBody{testErrorText}, r)
	}
}

func TestErrorBodyNoWrapsNotErrors(t *testing.T) {
	h := H(getFoo).With(ErrorBody{})
	s, r := h.Func(nil, nil, nil)
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
