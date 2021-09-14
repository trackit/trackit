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

package config

import (
	"testing"
)

const (
	envVarPrefix = "TRACKIT"
)

func TestIdentifierToEnvVarName(t *testing.T) {
	cases := [][2]string{
		{"HttpAddress", envVarPrefix + "_HTTP_ADDRESS"},
		{"SqlProtocol", envVarPrefix + "_SQL_PROTOCOL"},
		{"SqlAddress", envVarPrefix + "_SQL_ADDRESS"},
		{"HashDifficulty", envVarPrefix + "_HASH_DIFFICULTY"},
	}
	for _, c := range cases {
		if r := IdentifierToEnvVarName(c[0]); r != c[1] {
			t.Errorf("Testing with %s, should be %s, is %s.", c[0], c[1], r)
		}
	}
}
