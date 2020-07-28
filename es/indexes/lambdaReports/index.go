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

package lambdaReports

import "github.com/trackit/trackit/es/indexes/common"

const IndexSuffix = "lambda-reports"
const Type = "lambda-report"
const TemplateName = "lambda-reports"

type (
	// FunctionReport is saved in ES to have all the information of an Lambda function
	FunctionReport struct {
		common.ReportBase
		Function Function `json:"function"`
	}

	// FunctionBase contains basics information of an Lambda function
	FunctionBase struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		Version      string `json:"version"`
		LastModified string `json:"lastModified"`
		Runtime      string `json:"runtime"`
		Size         int64  `json:"size"`
		Memory       int64  `json:"memory"`
		Region       string `json:"region"`
	}

	// Function contains all the information of an Lambda function
	Function struct {
		FunctionBase
		Tags  []common.Tag `json:"tags"`
		Stats Stats        `json:"stats"`
	}

	// Stats contains statistics of an Lambda function
	Stats struct {
		Invocations Invocations `json:"invocations"`
		Duration    Duration    `json:"duration"`
	}

	Invocations struct {
		Total  float64 `json:"total"`
		Failed float64 `json:"failed"`
	}

	Duration struct {
		Average float64 `json:"average"`
		Maximum float64 `json:"maximum"`
	}
)
