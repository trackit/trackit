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
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/trackit/jsonlog"
)

// userKey is unique depending on user's aws' ID account and route's URL
func getUserKey(rdCache redisCache) string {
	key := fmt.Sprintf("%x-%x:", md5.Sum([]byte(rdCache.route)), md5.Sum([]byte(rdCache.args)))
	for _, val := range rdCache.awsAccount {
		key = fmt.Sprintf("%v%v:", key, val)
	}
	return key
}

func userHasCacheForService(rdCache redisCache, logger jsonlog.Logger) bool {
	userKey := getUserKey(rdCache)
	rtn, err := mainClient.Exists(userKey).Result()
	if err != nil {
		logger.Error("Unable to check if the key already exists in redis.", map[string] interface{} {
			"error":   err.Error(),
			"userKey": userKey,
		})
	}
	return err == nil && rtn > 0
}

func getUserCache(rdCache redisCache, logger jsonlog.Logger) interface{} {
	var cacheData interface{} = nil
	val, err := mainClient.Get(getUserKey(rdCache)).Result()
	if len(val) == 0 {
		logger.Error("No cache found on Redis for the despite it has been marked as valid", map[string]interface{}{
			"error": err})
		return nil
	}
	err = json.Unmarshal([]byte(val), &cacheData)
	if err != nil {
		logger.Error("Unable to unmarshal cache data for the key '%v'", map[string]interface{}{
			"error": err,
		})
		return nil
	}
	return cacheData
}

func createUserCache(rdCache redisCache, data interface{}, logger jsonlog.Logger) {
	if userHasCacheForService(rdCache, logger) {
		logger.Warning("The user has already a cache attributed for the current route.", map[string] interface{} {
			"userKey":   getUserKey(rdCache),
			"route":     rdCache.route,
			"routeArgs": rdCache.args,
			"accounts":  rdCache.awsAccount,
		})
		return
	}
	var err error
	rdCache.cacheContent, err = json.Marshal(data)
	if err != nil {
		logger.Error("Unable to marshal API content to create cache.", nil)
		return
	}
	userKey := getUserKey(rdCache)
	cmdStat := mainClient.Append(userKey, string(rdCache.cacheContent))
	if cmdStat.Err() == nil{
		mainClient.Expire(userKey, cacheExpireTime)
	} else {
		logger.Error("Unable to append content.", map[string] interface{} {
			"userKey": userKey,
			"error":   cmdStat.Err().Error(),
		})
	}
}

func deleteUserCache(rdCache redisCache, logger jsonlog.Logger) {
	rtn := mainClient.Del(getUserKey(rdCache))
	if rtn.Err() != nil {
		logger.Error("Unable to delete user's cache.", map[string] interface{} {
			"userKey": getUserKey(rdCache),
			"error":   rtn.Err().Error(),
		})
	}
}