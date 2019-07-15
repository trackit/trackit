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
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/config"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

// UsersCache is a struct to format a decorator that retrieve data from
// different route and cache it with redis which is an-memory data structure
// store, used as a database. Cache expire automatically after 24 hours.
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
		Addr:         config.RedisAddress,
		Password:     config.RedisPassword,
		DB:           1,
		IdleTimeout: -1 ,
	})
	_, err := mainClient.Ping().Result()
	if err != nil {
		jsonlog.DefaultLogger.Error("Unable to establish the connection to redis server", map[string] interface{} {
			"error": err.Error(),
		})
		log.Fatal("The API cannot connect to redis. Exiting...")
	}
	jsonlog.DefaultLogger.Info("Successfully connected to redis client", map[string] interface{} {
		"address": config.RedisAddress,
	})
}

// getFunc allows us to intercept the current data flow from the route and
// manipulate it to retrieve data or directly return data from the cache if
// there is one.
func (uc UsersCache) getFunc(hf routes.HandlerFunc) routes.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, args routes.Arguments) (int, interface{}) {
		logger := jsonlog.LoggerFromContextOrDefault(request.Context())
		_, userDataPresent := args[users.AuthenticatedUser].(users.User)
		_, dbAvailable := args[db.Transaction].(*sql.Tx)
		if !userDataPresent || !dbAvailable {
			logger.Error("Unable to retrieve user's or database's information while trying to get cache.", nil)
			writeHeaderCacheStatus(writer, cacheStatusError, "UNABLE-GET-BASICS-INFOS")
			return hf(writer, request, args)
		}
		rdCache, err := retrieveRouteInfos(request.URL.String(), args, logger)
		if err != nil {
			logger.Error("Error during cache initialization", map[string] interface{} {
				"error": err.Error(),
			})
			return hf(writer, request, args)
		}
		updateCacheByHeaderStatus(request, rdCache)
		if userHasCacheForService(rdCache, logger) {
			retrieveCache := getUserCache(rdCache, logger)
			if retrieveCache == nil {
				logger.Warning("Unable to retrieve cache, skipping it to avoid panic or error. The cache has been, also, deleted.", map[string] interface{} {
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
		if status == http.StatusOK {
			createUserCache(rdCache, routeData, logger)
			writeHeaderCacheStatus(writer, cacheStatusCreated)
		}
		return status, routeData
	}
}

func (uc UsersCache) Decorate(handler routes.Handler) routes.Handler {
	handler.Func = uc.getFunc(handler.Func)
	return handler
}