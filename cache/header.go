//   Copyright 2019 MSolution.IO
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

package cache

import (
	"net/http"

	"github.com/trackit/jsonlog"
)

// The different statuses below correspond to the cache status for a specified key.
// The communication is done through the "Cache-Status" header.
// When an error occurs, the cache status is set to "ERROR" and a key
// error message is indicated in the "Cache-Error" header.
const (
	cacheStatusCreated = iota
	cacheStatusUsed
	cacheStatusDelete
	cacheStatusError
	cacheStatusInvalid

	cacheHeaderAbsent  // Internal utilisation only

	cacheStatusKey    = "Cache-Status"
	cacheStatusErrKey = "Cache-Error"
)

// Keywords corresponding to each cache status.
// The front-end can only request the cache to be deleted with the keyword "DELETE".
var cacheHeaderStatus = map[int] string {
	cacheStatusCreated: "CREATED",
	cacheStatusUsed:    "USED",
	cacheStatusDelete:  "DELETE",
	cacheStatusError:   "ERROR",
}

func getHeaderCacheStatus(request *http.Request) int {
	status, exist := request.Header[cacheStatusKey]
	if !exist || len(status) == 0 {
		return cacheHeaderAbsent
	}
	for i := range cacheHeaderStatus {
		if cacheHeaderStatus[i] == status[0] {
			return i
		}
	}
	return cacheStatusInvalid
}

func writeHeaderCacheStatus(writer http.ResponseWriter, newStatus int, errorOccur ...string) {
	if keyContent := writer.Header().Get(cacheStatusKey); len(keyContent) == 0 {
		writer.Header().Add(cacheStatusKey, cacheHeaderStatus[newStatus])
	} else {
		writer.Header().Set(cacheStatusKey, cacheHeaderStatus[newStatus])
	}
	if newStatus == cacheStatusError {
		if errMsg := writer.Header().Get(cacheStatusErrKey); len(errMsg) == 0 {
			writer.Header().Add(cacheStatusErrKey, errorOccur[0])
		} else {
			writer.Header().Set(cacheStatusErrKey, errorOccur[0])
		}
	}
}

func updateCacheByHeaderStatus(request *http.Request, rdCache redisCache) {
	status := getHeaderCacheStatus(request)
	if status == cacheStatusDelete {
		deleteUserCache(rdCache, jsonlog.LoggerFromContextOrDefault(request.Context()))
	}
}