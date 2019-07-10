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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/config"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/models"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

const cacheExpireTime = 24 * time.Hour

// UsersCache is a struct to format a decorator that retrieve data from
// different route and cache it with redis which is an-memory data structure
// store, used as a database. Cache expire automatically after 24 hours.
type UsersCache struct {}

var mainClient *redis.Client

func init() {
	mainClient = redis.NewClient(&redis.Options{
		Addr:        config.RedisAddress,
		Password:    config.RedisPassword,
		DB:          1,
		IdleTimeout: 10,
	})
	_, err := mainClient.Ping().Result()
	if err != nil {
		log.Fatalf("Cannot connect to redis. Got error: '%v'", err.Error())
	}
	log.Printf("Successfully connected to client redis at adress '%v'.\n", config.RedisAddress)
}

// userKey is unique depending on user's e-mail address and route's URL
func getUserKey(identity, service string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(identity + service))
}

func userHasCacheForService(identity, service string) bool {
	rtn, err := mainClient.Exists(getUserKey(identity, service)).Result()
	return rtn > 0 && err == nil
}

func getAwsAccIdFromRaw(args routes.Arguments) string {
	userId := args[users.AuthenticatedUser].(users.User)
	context := args[db.Transaction].(*sql.Tx)
	awsAcc, err := models.AwsAccountByID(context, userId.Id)
	if err != nil {
		return "unknown"
	} else {
		return awsAcc.AwsIdentity
	}
}

func getUserCache(identity, service string, logger jsonlog.Logger) interface{} {
	var cacheData interface{} = nil
	rtn, err := mainClient.Get(getUserKey(identity, service)).Bytes()
	if err != nil {
		_ = logger.Error(fmt.Sprintf("No cache found on Redis for the '%v' despite it has marked as valid\n", identity), err)
		return nil
	} else {
		err := json.Unmarshal(rtn, &cacheData)
		if err != nil {
			_ = logger.Error(fmt.Sprintf("Unable to unmarshal cache data for the key '%v'\n", identity), err)
			return nil
		}
		return cacheData
	}
}

func createUserCache(identity, service string, data interface{}, logger jsonlog.Logger) error {
	if userHasCacheForService(identity, service) {
		_ = logger.Warning(fmt.Sprintf("User %v has already a cache attributed for service '%v'", identity, service), nil)
		return nil
	}
	content, newErr := json.Marshal(data)
	if newErr != nil {
		_ = logger.Error("Unable to marshal API content to create cache.", nil)
		return newErr
	}
	userKey := getUserKey(identity, service)
	err := mainClient.Append(userKey, string(content))
	mainClient.Expire(userKey, cacheExpireTime)
	return err.Err()
}

func deleteUserCache(identity, service string) {
	mainClient.Del(getUserKey(identity, service))
}

// getFunc allows us to intercept the current data flow from the route and
// manipulate it to retrieve data or directly return data from the cache if
// there is one.
func (_ UsersCache) getFunc(hf routes.HandlerFunc) routes.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, args routes.Arguments) (int, interface{}) {
		logger := jsonlog.LoggerFromContextOrDefault(request.Context())
		_, userDataPresent := args[users.AuthenticatedUser]
		_, dbAvailable := args[db.Transaction]
		if !userDataPresent || !dbAvailable {
			_ = logger.Error("Unable to retrieve user's or database's information while trying to get cache.", nil)
			return hf(writer, request, args)
		}
		service := request.URL.String()
		if identity := getAwsAccIdFromRaw(args); userHasCacheForService(identity, service) {
			retrieveCache := getUserCache(identity, service, jsonlog.LoggerFromContextOrDefault(request.Context()))
			if retrieveCache == nil {
				_ = logger.Warning(fmt.Sprintf("Unable to retrieve cache for key %v, skipping it to avoid panic or error. The cache has been, also, deleted.\n", getUserKey(identity, service)), nil)
				deleteUserCache(identity, service)
			} else {
				return http.StatusOK, retrieveCache
			}
		}
		status, routeData := hf(writer, request, args)
		if status == http.StatusOK {
			_ = createUserCache(getAwsAccIdFromRaw(args), service, routeData, jsonlog.LoggerFromContextOrDefault(request.Context()))
		}
		return status, routeData
	}
}

func (uc UsersCache) Decorate(handler routes.Handler) routes.Handler {
	handler.Func = uc.getFunc(handler.Func)
	return handler
}