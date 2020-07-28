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

package ebsReports

import (
	"time"

	"github.com/trackit/trackit/es/indexes/common"
)

const IndexSuffix = "ebs-reports"
const Type = "ebs-report"
const TemplateName = "ebs-reports"

type (
	// SnapshotReport is saved in ES to have all the information of an EBS snapshot
	SnapshotReport struct {
		common.ReportBase
		Snapshot Snapshot `json:"snapshot"`
	}

	// SnapshotBase contains basics information of an EBS snapshot
	SnapshotBase struct {
		Id          string    `json:"id"`
		Description string    `json:"description"`
		State       string    `json:"state"`
		Encrypted   bool      `json:"encrypted"`
		StartTime   time.Time `json:"startTime"`
		Region      string    `json:"region"`
	}

	// Snapshot contains all the information of an EBS snapshot
	Snapshot struct {
		SnapshotBase
		Tags   []common.Tag `json:"tags"`
		Volume Volume       `json:"volume"`
		Cost   float64      `json:"cost"`
	}

	// Volume contains information about an EBS volume
	Volume struct {
		Id   string `json:"id"`
		Size int64  `json:"size"`
	}
)
