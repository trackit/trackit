//   Copyright 2018 MSolution.IO
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

package plugins_utils

// StatusPercentSteps defines the steps used to compute the status
// Example: StatusPercentSteps{25, 75}
// A percentage between 0 and 24 will have a red status
// A percentage between 25 and 74 will have a orange status
// A percentage between 75 and 100 will have a green status
// A percentage < 0 or > 100 will have a red status
type StatusPercentSteps struct {
	MinOrange int
	MinGreen  int
}

// GetStatus returns the status based on the parameters and on the configured steps
func (steps StatusPercentSteps) GetStatus(checked, passed int) string {
	if checked == 0 {
		return "green"
	}
	percentSuccess := int(float64(passed) / float64(checked) * 100.0)
	if percentSuccess >= steps.MinGreen && percentSuccess <= 100 {
		return "green"
	} else if percentSuccess >= steps.MinOrange && percentSuccess < steps.MinGreen {
		return "orange"
	}
	return "red"
}
