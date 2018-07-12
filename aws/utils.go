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

package aws

import (
	"fmt"
)

// ValidateAwsAccounts will validate a slice of int passed to it.
// It checks that they are 12 digit numbers
func ValidateAwsAccounts(awsAccounts []string) error {
	for _, account := range awsAccounts {
		if len(account) != 12 {
			return fmt.Errorf("invalid account format : %s", account)
		}
	}
	return nil
}
