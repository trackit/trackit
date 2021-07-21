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

import (
	"fmt"
	"time"

	utils "github.com/trackit/trackit/aws/usageReports"
)

const urlFormat = "https://console.aws.amazon.com/cloudformation/home?region=%s#/stacks/stackinfo?stackId=%s"
const stackId = "aws:cloudformation:stack-id"

func (doc *TaggingReportDocument) GenCloudFormationUrl() {
	for _, tag := range doc.Tags {
		if tag.Key == stackId && tag.Value != "" {
			doc.CloudFormationURL = fmt.Sprintf(urlFormat, doc.Region, tag.Value)
		}
	}
}

// TaggingReportDocument is an entry in ES' tagging index
type TaggingReportDocument struct {
	Account           string      `json:"account"`
	ReportDate        time.Time   `json:"reportDate"`
	ResourceID        string      `json:"resourceId"`
	ResourceType      string      `json:"resourceType"`
	Region            string      `json:"region"`
	URL               string      `json:"url"`
	Tags              []utils.Tag `json:"tags"`
	CloudFormationURL string      `json:"cloudFormationUrl,omitempty"`
}
