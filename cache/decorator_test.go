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
	"fmt"
	"testing"
)

const (
	testIdentity = "team@trackit.io"
	testRoute    = "/lambda"
	testArgs     = "date=2019-07-02"
	testAwsAcc   = "394125495069"

	testKey      = "KEY"
	testContent  = "This is the content I like"
)

func TestGetUserKey(test *testing.T) {
	result := getUserKey(redisCache{
		route:      testRoute,
		args:       testArgs,
		awsAccount: []string {testAwsAcc},
	})
	excepted := fmt.Sprintf("%x-%x:%v:", md5.Sum([]byte(testRoute)), md5.Sum([]byte(testArgs)), testAwsAcc)
	if result != excepted {
		test.Errorf("Execepted '%v' but got '%v'", excepted, result)
	}
}