//   Copyright 2020 MSolution.IO
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

package utils

import "strconv"

// GetRegionForURL removes the last character of an AWS region if it's not a number
func GetRegionForURL(region string) string {
	if len(region) <= 0 {
		return region
	}

	_, err := strconv.ParseInt(region[len(region)-1:], 10, 32)
	if err != nil {
		region = region[:len(region)-1]
	}

	return region
}
