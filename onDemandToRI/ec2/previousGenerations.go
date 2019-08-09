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

package onDemandToRiEc2

var (
	// PreviousToCurrentGeneration maps previous generation instance type to current
	// generation instance types
	PreviousToCurrentGeneration = map[string]string{
		"t1.micro":    "t3.nano",
		"m3.medium":   "m5.large",
		"m3.large":    "m5.large",
		"m3.xlarge":   "m5.xlarge",
		"m3.2xlarge":  "m5.2xlarge",
		"m1.small":    "t2.small",
		"m1.medium":   "t2.medium",
		"m1.large":    "m5.large",
		"m1.xlarge":   "m5.xlarge",
		"m2.xlarge":   "r5.large",
		"m2.2xlarge":  "r5.xlarge",
		"m2.4xlarge":  "r5.2xlarge",
		"r3.large":    "r5.large",
		"r3.xlarge":   "r5.xlarge",
		"r3.2xlarge":  "r5.2xlarge",
		"r3.4xlarge":  "r5.4xlarge",
		"r3.8xlarge":  "r4.8xlarge",
		"cr1.8xlarge": "r4.8xlarge",
		"c3.large":    "c5.large",
		"c3.xlarge":   "c5.xlarge",
		"c3.2xlarge":  "c5.2xlarge",
		"c3.4xlarge":  "c5.4xlarge",
		"c3.8xlarge":  "c5.9xlarge",
		"c1.medium":   "c5.large",
		"c1.xlarge":   "c5.2xlarge",
		"cc2.8xlarge": "c5.9xlarge",
		"i2.xlarge":   "i3.xlarge",
		"i2.2xlarge":  "i3.2xlarge",
		"i2.4xlarge":  "i3.4xlarge",
		"i2.8xlarge":  "i3.8xlarge",
		"hs1.8xlarge": "i3.8xlarge",
		"g2.2xlarge":  "g3.4xlarge",
		"g2.8xlarge":  "g3.8xlarge",
	}
)
