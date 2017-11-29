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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/trackit/jsonlog"
)

const panicPayload = "panic!"

func panickingHandler(_ *http.Request, _ Arguments) (int, interface{}) {
	panic(panicPayload)
}

func requestWithLogger(r *http.Request, l jsonlog.Logger) *http.Request {
	ctx := r.Context()
	ctx = jsonlog.ContextWithLogger(ctx, l)
	return r.WithContext(ctx)
}

func TestPanicError(t *testing.T) {
	h := H(panickingHandler).With(PanicAsError{})
	logger := jsonlog.DefaultLogger.WithWriter(ioutil.Discard)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request = requestWithLogger(request, logger)
	s, r := h.Func(nil, request, nil)
	if s != http.StatusInternalServerError {
		t.Errorf("Status code should be %d, is %d instead.", http.StatusOK, s)
	}
	if rt, ok := r.(error); ok {
		if rt != panicResponse {
			t.Errorf("Response should be '%s', is '%s' instead.", panicResponse, rt.Error())
		}
	} else {
		t.Errorf("Response should be %[1]T %#[1]v, is %[2]T %#[2]v instead.", panicResponse, r)
	}
}
