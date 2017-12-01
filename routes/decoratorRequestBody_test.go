package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type bodyEchoTest struct {
	Foo string `json:"foo" req:"nonzero"`
	Bar int    `json:"bar"`
}

func bodyEchoHandler(r *http.Request, a Arguments) (int, interface{}) {
	var body bodyEchoTest
	MustRequestBody(a, &body)
	return http.StatusOK, body
}

func TestEchoHandler(t *testing.T) {
	example := bodyEchoTest{"test", 42}
	arguments := make(Arguments)
	arguments[argumentKeyBody] = reflect.ValueOf(example)
	status, response := bodyEchoHandler(nil, arguments)
	if status != http.StatusOK {
		t.Errorf("Status code should be %d, is %d.", http.StatusOK, status)
	}
	if response != example {
		t.Errorf("Response should be %v, is %v.", example, response)
	}
}

func TestSuccessfulValidation(t *testing.T) {
	example := bodyEchoTest{"test", 42}
	handler := H(bodyEchoHandler).With(
		RequestBody{example},
	)
	bs, _ := json.Marshal(example)
	requestBody := bytes.NewBuffer(bs)
	request := httptest.NewRequest(http.MethodGet, "/", requestBody)
	status, response := handler.Func(nil, request, make(Arguments))
	if status != http.StatusOK {
		t.Errorf("Status code should be %d, is %d.", http.StatusOK, status)
	}
	if response != example {
		t.Errorf("Response should be %v, is %v.", example, response)
	}
}

func TestFailedValidation(t *testing.T) {
	const expectedError = "foo: value is zero"
	example := bodyEchoTest{"", 42}
	handler := H(bodyEchoHandler).With(
		RequestBody{example},
	)
	bs, _ := json.Marshal(example)
	requestBody := bytes.NewBuffer(bs)
	request := httptest.NewRequest(http.MethodGet, "/", requestBody)
	status, response := handler.Func(nil, request, make(Arguments))
	if status != http.StatusBadRequest {
		t.Errorf("Status code should be %d, is %d.", http.StatusBadRequest, status)
	}
	if errorResponse, ok := response.(error); !ok {
		t.Errorf("Response should be error %q, is %v.", expectedError, response)
	} else if errorResponse.Error() != expectedError {
		t.Errorf("Response should be error %q, is error %q.", expectedError, errorResponse.Error())
	}
}

const expectedDoc = `{
	"summary": "",
	"components": {
		"input:body:example": {
			"summary": "input body example",
			"description": "{\n\t\"foo\": \"test\",\n\t\"bar\": 42\n}"
		}
	}
}`

func TestDocumentation(t *testing.T) {
	example := bodyEchoTest{"test", 42}
	handler := H(bodyEchoHandler).With(
		RequestBody{example},
	)
	doc, _ := json.MarshalIndent(handler.Documentation, "", "\t")
	sdoc := string(doc)
	if sdoc != expectedDoc {
		t.Errorf("Documentation should be %q, is %q.", expectedDoc, sdoc)
	}
}
