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
	"encoding/json"

	"github.com/trackit/jsonlog"
)

func userHasCacheForService(rdCache redisCache, logger jsonlog.Logger) bool {
	rtn, err := mainClient.Exists(rdCache.key).Result()
	if err != nil {
		logger.Error("Unable to check if the key already exists in redis.", map[string] interface{} {
			"error":   err.Error(),
			"userKey": rdCache.key,
		})
	}
	return err == nil && rtn > 0
}

func getUserCache(rdCache redisCache, logger jsonlog.Logger) interface{} {
	var cacheData interface{} = nil
	val, err := mainClient.Get(rdCache.key).Result()
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
			"userKey":   rdCache.key,
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
	cmdStat := mainClient.Append(rdCache.key, string(rdCache.cacheContent))
	if cmdStat.Err() == nil{
		mainClient.Expire(rdCache.key, cacheExpireTime)
	} else {
		logger.Error("Unable to append content.", map[string] interface{} {
			"userKey": rdCache.key,
			"error":   cmdStat.Err().Error(),
		})
	}
}

func deleteUserCache(rdCache redisCache, logger jsonlog.Logger) {
	rtn := mainClient.Del(rdCache.key)
	if rtn.Err() != nil {
		logger.Error("Unable to delete user's cache.", map[string] interface{} {
			"userKey": rdCache.key,
			"error":   rtn.Err().Error(),
		})
	}
}