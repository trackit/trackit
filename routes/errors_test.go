package routes

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

const errorText = "Hello world!"
const successText = "Good!"

func errorHandler(r *http.Request, a Arguments) (int, interface{}) {
	return 500, errors.New(errorText)
}

func successHandler(r *http.Request, a Arguments) (int, interface{}) {
	return 200, successText
}

func TestErrorWrapped(t *testing.T) {
	h := ApplyDecorators(
		baseIntermediate(errorHandler),
		WithErrorBody{},
	)
	request := httptest.NewRequest("GET", "/error", nil)
	response := httptest.NewRecorder()
	status, body := h(response, request, nil)
	if status != 500 {
		t.Errorf("Expected status 500, got %d.", status)
	}
	if errorBody, ok := body.(ErrorBody); !ok {
		t.Errorf("Expected ErrorBody, got %T %#v.", body, body)
	} else if errorBody.Error != errorText {
		t.Errorf("Expected error to be '%s', got '%s'.", errorText, errorBody.Error)
	}
}

func TestSuccessForwarded(t *testing.T) {
	h := ApplyDecorators(
		baseIntermediate(successHandler),
		WithErrorBody{},
	)
	request := httptest.NewRequest("GET", "/success", nil)
	response := httptest.NewRecorder()
	status, body := h(response, request, nil)
	if status != 200 {
		t.Errorf("Expected status 200, got %d.", status)
	}
	if successBody, ok := body.(string); !ok {
		t.Errorf("Expected string, got %T %#v.", body, body)
	} else if successBody != successText {
		t.Errorf("Expected success to be '%s', got '%s'.", successText, successBody)
	}
}
