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
	"context"
	"fmt"
	"net/http"

	"github.com/trackit/jsonlog"
)

// PanicAsError decorates handlers to recover from panics: it uses the panic's
// payload and wraps it into an error. It also logs the panic.
type PanicAsError struct{}

func (d PanicAsError) Decorate(h Handler) Handler {
	h.Func = d.getFunc(h.Func)
	return h
}

// getFunc builds the handler function for PanicAsError.Decorate.
func (_ PanicAsError) getFunc(hf HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (status int, output interface{}) {
		defer func() {
			if rc := recover(); rc != nil {
				status, output = handlePanic(r.Context(), rc)
			}
		}()
		return hf(w, r, a)
	}
}

const panicResponse = constError("There was an internal server error.")

func handlePanic(ctx context.Context, r interface{}) (int, interface{}) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	l.Error("Route handler panicked.", fmt.Sprintf("%v", r))
	return http.StatusInternalServerError, panicResponse
}
