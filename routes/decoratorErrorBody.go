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
)

// ErrorBody is a decorator for an HTTP handler. If that handler returns an
// error, it will wrap it in a structure so that it can correctly be returned
// to the user as JSON.
type ErrorBody struct{}

type errorBody struct {
	Error string `json:"error"`
}

func (d ErrorBody) Decorate(h Handler) Handler {
	h.Func = d.getFunc(h.Func)
	return h
}

// getFunc builds the handler function for ErrorBody.Decorate.
func (_ ErrorBody) getFunc(hf HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		status, output := hf(w, r, a)
		if err, ok := output.(error); ok {
			return status, errorBody{err.Error()}
		} else {
			return status, output
		}
	}
}
