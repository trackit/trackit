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

// Package cache implements caching of results in Redis
package cache

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"

	"github.com/trackit/jsonlog"
)

func getTotalRedisKeys() (totalKeys int64, err error) {
	const subStrLen = len("keys=")
	// The format of retrieving info of the space "keyspace" is has the format below:
	// "dbW:keys=X,expires=Y,avg_ttl=Z"
	// We want retrieve te total of keys, which is X in the format.
	keySpace := mainClient.Info("keyspace")
	idx := strings.Index(keySpace.Val(), "keys=") // Idx is equal to the index of letter 'k' from "keys="
	if idx == -1 {
		return
	}
	// In "keys=X,expires=Y,avg_ttl=Z", we look for the index of the first comma
	// We have to subtract the length of "keys=" (which is subStrLen) to know the number length.
	comma := strings.IndexByte(keySpace.Val()[idx:], ',') - subStrLen
	if comma == -1 {
		return
	}
	// Start at idx + subStrLen means the first number character and
	// the previous number + comma correspond to the last number character
	starterIdx := idx + subStrLen
	endIdx := starterIdx + comma
	return strconv.ParseInt(keySpace.Val()[starterIdx:endIdx], 10, 64)
}

// RemoveMatchingCache removes all cache related to the format ROUTE-...-KEY-
// It's important to note that AWS identities and routes validity isn't checked.
func RemoveMatchingCache(routes []string, awsAccounts []string, logger jsonlog.Logger) (err error) {
	var totalKeys int64
	totalKeys, err = getTotalRedisKeys()
	if err != nil {
		logger.Error("Unable to get the total redis keys.", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if totalKeys == 0 {
		return nil
	}
	const matchPattern = "%x-*-%v-*"
	for _, route := range routes {
		routeKey := md5.Sum([]byte(route))
		for _, awsAcc := range awsAccounts {
			obsoleteKeys := mainClient.Scan(0, fmt.Sprintf(matchPattern, routeKey, awsAcc), totalKeys)
			if obsoleteKeys.Err() != nil {
				logger.Error("Unable to scan redis DB to retrieve matching keys.", map[string]interface{}{
					"error":       obsoleteKeys.Err().Error(),
					"awsIdentity": awsAcc,
					"route":       route,
					"keyFormat":   fmt.Sprintf(matchPattern, routeKey, awsAcc),
				})
				continue
			}
			cacheKeys, _ := obsoleteKeys.Val()
			for _, keyValue := range cacheKeys {
				rtn := mainClient.Del(keyValue)
				if rtn.Err() != nil {
					logger.Warning("Unable to delete cache for a specific key from a route.", map[string]interface{}{
						"error": rtn.Err().Error(),
						"route": route,
						"key":   keyValue,
					})
				}
			}
		}
	}
	return
}
