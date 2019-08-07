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
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/go-redis/redis"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// UsersCache is a struct to format a decorator that retrieve data from
// different route and cache it with redis. Cache expire automatically
// after 24 hours.
type UsersCache struct {
}

type redisCache struct {
	route        string
	args         string
	awsAccount   []string
	key          string
	cacheContent []byte
}

const cacheExpireTime = 24 * time.Hour

var mainClient *redis.Client

func init() {
	mainClient = redis.NewClient(&redis.Options{
		Addr:        config.RedisAddress,
		Password:    config.RedisPassword,
		DB:          1,
		IdleTimeout: -1,
	})
	_, err := mainClient.Ping().Result()
	if err != nil {
		jsonlog.Error("Unable to establish the connection to redis server", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}
	jsonlog.Info("Successfully connected to redis client", map[string]interface{}{
		"address": config.RedisAddress,
	})
}

// getFunc allows us to intercept the current data flow from the route and
// manipulate it to retrieve data or directly return data from the cache if
// there is one.
func (uc UsersCache) getFunc(hf routes.HandlerFunc) routes.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, args routes.Arguments) (int, interface{}) {
		logger := jsonlog.LoggerFromContextOrDefault(request.Context())
		if _, userDataPresent := args[users.AuthenticatedUser].(users.User); !userDataPresent {
			logger.Error("Unable to retrieve user's information while trying to get cache.", nil)
			writeHeaderCacheStatus(writer, cacheStatusError, "UNABLE-GET-BASICS-INFOS-USER")
			return hf(writer, request, args)
		} else if _, dbAvailable := args[db.Transaction].(*sql.Tx); !dbAvailable {
			logger.Error("Unable to retrieve database's information while trying to get cache.", nil)
			writeHeaderCacheStatus(writer, cacheStatusError, "UNABLE-GET-BASICS-INFOS-DB")
			return hf(writer, request, args)
		}
		rdCache, err := initialiseCacheInfos(request.URL.String(), args, logger)
		if err != nil {
			logger.Error("Error during cache initialization", map[string]interface{}{
				"error": err.Error(),
			})
			writeHeaderCacheStatus(writer, cacheStatusError, "ERROR-INITIALIZATION")
			return hf(writer, request, args)
		}
		updateCacheByHeaderStatus(request, rdCache)
		if userHasCacheForService(rdCache, logger) {
			retrieveCache := getUserCache(rdCache, logger)
			if retrieveCache == nil {
				logger.Warning("Unable to retrieve cache, skipping it to avoid panic or error. The cache has been deleted.", map[string]interface{}{
					"userKey": rdCache.key,
					"route":   rdCache.route,
				})
				deleteUserCache(rdCache, logger)
			} else {
				writeHeaderCacheStatus(writer, cacheStatusUsed)
				return http.StatusOK, retrieveCache
			}
		}
		status, routeData := hf(writer, request, args)
		if status == http.StatusOK && isValidResponse(routeData) {
			createUserCache(rdCache, routeData, logger)
			writeHeaderCacheStatus(writer, cacheStatusCreated)
		}
		return status, routeData
	}
}

func isValidResponse(data interface{}) bool {
	exp, err := regexp.Compile(`{.*}|\[.*\]`)
	if err != nil {
		return false
	}
	return exp.MatchString(fmt.Sprintf("%v", data))
}

func (uc UsersCache) Decorate(handler routes.Handler) routes.Handler {
	handler.Func = uc.getFunc(handler.Func)
	return handler
}
