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
