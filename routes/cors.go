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

// WithCors is a route decorator which adds some CORS headers to all responses
// in a static way. This is not in my knowledge correct behavior but it will do
// for now.
type WithCors struct {
	AllowOrigin      []string
	AllowHeaders     []string
	AllowMethods     []string
	AllowCredentials bool
}

func (d WithCors) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		w.Header()["Access-Control-Allow-Origin"] = []string{strings.Join(d.AllowOrigin, ",")}
		w.Header()["Access-Control-Allow-Methods"] = []string{strings.Join(d.AllowMethods, ",")}
		w.Header()["Access-Control-Allow-Headers"] = []string{strings.Join(d.AllowHeaders, ",")}
		w.Header()["Access-Control-Allow-Credentials"] = []string{fmt.Sprintf("%v", d.AllowCredentials)}
		return h(w, r, a)
	}
}
