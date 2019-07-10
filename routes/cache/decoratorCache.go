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
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/config"
	"github.com/trackit/trackit-server/users"
)

const cacheExpireTime = 2 * time.Second

type clientCache struct {
	service map[string] interface{}
}

type UsersCache struct {
	ServiceName string
}

var (
	mainClient *redis.Client
	usersCache map[string] clientCache
	errNonexistentCache = errors.New("trying to retrieve a nonexistent user cache")
	errNonexistentService = errors.New("trying to retrieve a nonexistent service")
)

func init() {
	mainClient = redis.NewClient(&redis.Options{
		Addr:        config.RedisAddress,
		Password:    config.RedisPassword,
		DB:          1,
		IdleTimeout: 10,
	})
	_, err := mainClient.Ping().Result()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Client ('%v') has been successfully connected: %+v\n", config.RedisAddress, mainClient)
}

func (uc UsersCache) Decorate(handler Handler) Handler {
	handler.Func = uc.getFunc(handler.Func)
	return handler
}

func getUserCache(identity, service string) interface{} {
	if _, userExist := usersCache[identity]; userExist {
		if _, serviceExist := usersCache[identity].service[service]; serviceExist {
			return usersCache[identity].service[service]
		}
		return errNonexistentService
	}
	return errNonexistentCache
}

func (_ UsersCache) getFunc(hf HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (int, interface{}) {
		if service := r.URL.String(); UserHasCacheForService(a[users.AuthenticatedUser].(string), service) {
			return http.StatusOK, getUserCache(a[users.AuthenticatedUser].(string), service)
		}
		return hf(w, r, a)
	}
}

func UserHasCacheForService(identity, service string) bool {
	_, userExist := usersCache[identity]
	if userExist {
		_, serviceExist := usersCache[identity].service[service]
		return serviceExist
	}
	return userExist
}

// awsAcc *models.AwsAccount
func CreateUserCache(identity, service string, data interface{}) error {
	if UserHasCacheForService(identity, service) {
		fmt.Printf("User %v has already a cache attributed for service '%v'", identity, service)
		return nil
	}
	mainClient.Append("key", "VAL")
	mainClient.Expire("key", cacheExpireTime)
	if rtn := mainClient.Exists("key"); rtn.Val() > 0 {
		fmt.Printf("Key exists\n")
	}
	defer func() {
		fmt.Printf("In defered funct\n")
		if mainClient.Exists("key").Val() > 0 {
			fmt.Printf("Key exists\n")
		} else {
			fmt.Print("Key doesn't exist anymore\n")
		}
	}()
	time.Sleep(4)
}