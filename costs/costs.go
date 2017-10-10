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

package costs

import (
	"net/http"

	"github.com/trackit/jsonlog"
)

// HandleRequest is a dummy request handler function. It does nothing except
// some logging and returns static data.
func HandleRequest(response http.ResponseWriter, request *http.Request) {
	logger := jsonlog.LoggerFromContextOrDefault(request.Context())
	logger.Debug("Request headers.", request.Header)
	response.WriteHeader(200)
	response.Write([]byte("Costs."))
}
