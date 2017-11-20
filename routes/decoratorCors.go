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
	"fmt"
	"net/http"
	"strings"
)

// Cors is a decorator which adds support for CORS to the handler.
type Cors struct {
	AllowOrigin      []string
	AllowHeaders     []string
	AllowCredentials bool
}

func (c Cors) Decorate(h Handler) Handler {
	mm := h.methods
	ms := make([]string, 0, len(mm)+1)
	ms = append(ms, http.MethodOptions)
	for m, v := range mm {
		if v {
			ms = append(ms, m)
		}
	}
	acaMethods := []string{strings.Join(ms, ",")}
	acaOrigin := []string{strings.Join(c.AllowOrigin, ",")}
	acaHeaders := []string{strings.Join(c.AllowHeaders, ",")}
	acaCredentials := []string{fmt.Sprintf("%t", c.AllowCredentials)}
	h.Func = c.getFunc(h.Func, h.Documentation, acaMethods, acaOrigin, acaHeaders, acaCredentials)
	return h
}

// getFunc builds a handler function for Cors.Decorate.
func (_ Cors) getFunc(
	hf HandlerFunc,
	hd HandlerDocumentation,
	acaMethods []string,
	acaOrigin []string,
	acaHeaders []string,
	acaCredentials []string,
) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		w.Header()["Access-Control-Allow-Methods"] = acaMethods
		w.Header()["Access-Control-Allow-Origin"] = acaOrigin
		w.Header()["Access-Control-Allow-Headers"] = acaHeaders
		w.Header()["Access-Control-Allow-Credentials"] = acaCredentials
		if r.Method == http.MethodOptions {
			return http.StatusOK, hd
		} else {
			return hf(w, r, a)
		}
	}
}
