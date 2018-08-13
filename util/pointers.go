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

package util

func SafeStringFromPtr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func SafeInt64FromPtr(nb *int64) int64 {
	if nb == nil {
		return int64(0)
	}
	return *nb
}

func SafeBoolFromPtr(val *bool) bool {
	if val == nil {
		return false
	}
	return *val
}
