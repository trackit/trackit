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

package reports

import (
	"context"
	"database/sql"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/trackit/trackit/aws"
)

type module struct {
	Name          string
	SheetName     string
	ErrorName     string
	GenerateSheet func(context.Context, []aws.AwsAccount, time.Time, *sql.Tx, *excelize.File) error
}

var modules = []module{
	costVariationLastMonth,
	costVariationLast6Months,
	ec2UsageReportModule,
	rdsUsageReportModule,
	esUsageReportModule,
	lambdaUsageReportModule,
	elastiCacheUsageReportModule,
	s3CostReportModule,
	ebsUsageReportModule,
	instanceCountUsageReportModule,
}
