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
	"github.com/go-redis/redis"
	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit-server/users"
)

type UserCache struct {
	Client* redis.Client
	User 	users.User
}

//var UsersCache = make(map[string] *UserCache)

const cacheExpireTime = 360

func CreateUserCache(user users.User, service string, data interface{}) error {
	//if _, exists := UsersCache[user.Email]; exists {
	//	_ = jsonlog.Error(fmt.Sprintf("User %v has already a cache for service %v"), nil)
	//}
	//UsersCache[user.Email].Client = redis.NewClient(&redis.Options{
	//	Addr: "127.0.0.1:6379",
	//	Password: "",
	//	DB: 0})

	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if client == nil {
		_ = jsonlog.Error(fmt.Sprintf("Unable to make a new client for %v's cache", user.Email), nil)
		return nil
	}
	//client.Set(service, data, cacheExpireTime)
	pong, err := client.Ping().Result()
	fmt.Printf("Pong: '%v' & err: '%v'\n", pong, err)
	return nil
}