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

package esReports

import "github.com/trackit/trackit/es/indexes/common"

const IndexSuffix = "es-reports"
const Type = "es-report"
const TemplateName = "es-reports"

type (
	// DomainReport represents the report with all the information for ES domains.
	DomainReport struct {
		common.ReportBase
		Domain Domain `json:"domain"`
	}

	// DomainBase contains basics information of an ES domain
	DomainBase struct {
		Arn               string `json:"arn"`
		Region            string `json:"region"`
		DomainID          string `json:"domainId"`
		DomainName        string `json:"domainName"`
		InstanceType      string `json:"instanceType"`
		InstanceCount     int64  `json:"instanceCount"`
		TotalStorageSpace int64  `json:"totalStorageSpace"`
	}

	// Domain contains all information of an ES domain that will be save in ES
	Domain struct {
		DomainBase
		Tags  []common.Tag       `json:"tags"`
		Costs map[string]float64 `json:"costs"`
		Stats Stats              `json:"stats"`
	}

	// Stats contains statistics of a domain get on CloudWatch
	Stats struct {
		Cpu               Cpu               `json:"cpu"`
		FreeSpace         float64           `json:"freeSpace"`
		JVMMemoryPressure JVMMemoryPressure `json:"JVMMemoryPressure"`
	}

	// Cpu contains cpu statistics of a domain
	Cpu struct {
		Average float64 `json:"average"`
		Peak    float64 `json:"peak"`
	}

	// JVMMemoryPressure contains JVMMemoryPressure statistics of a domain
	JVMMemoryPressure struct {
		Average float64 `json:"average"`
		Peak    float64 `json:"peak"`
	}
)
