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

package instanceCountReports

import (
	"time"

	"github.com/trackit/trackit/es/indexes/common"
)

const IndexSuffix = "instancecount-reports"
const Type = "instancecount-report"
const TemplateName = "instancecount-reports"

type (
	// InstanceCount is saved in ES to have all the information of an InstanceCount
	InstanceCountReport struct {
		common.ReportBase
		InstanceCount InstanceCount `json:"instanceCount"`
	}

	// InstanceCount contains all the information of an InstanceCount
	InstanceCount struct {
		Type   string               `json:"instanceType"`
		Region string               `json:"region"`
		Hours  []InstanceCountHours `json:"hours"`
	}

	InstanceCountHours struct {
		Hour  time.Time `json:"hour"`
		Count float64   `json:"count"`
	}
)
