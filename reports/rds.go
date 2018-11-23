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

package reports

import (
	"context"
	"database/sql"
	"github.com/trackit/jsonlog"
	"strconv"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/rds"
	"github.com/trackit/trackit-server/users"
)

var rdsInstanceFormat = []string{
	"Name",
	"Type",
	"Region",
	"Cost",
	"Engine",
	"Multi A-Z",
	"Storage (GigaBytes)",
	"Storage - Available (Average) (Bytes)",
	"Storage - Available (Minimum) (Bytes)",
	"Storage - Available (Maximum) (Bytes)",
	"CPU Average (Percentage)",
	"CPU Peak (Percentage)",
}

func formatRdsInstance(instance rds.Instance) []string {
	var cost float64

	for _, value := range instance.Costs {
		cost += value
	}

	return []string{
		instance.DBInstanceIdentifier,
		instance.DBInstanceClass,
		instance.AvailabilityZone,
		strconv.FormatFloat(cost, 'f', -1, 64),
		instance.Engine,
		strconv.FormatBool(instance.MultiAZ),
		strconv.FormatInt(instance.AllocatedStorage, 10),
		strconv.FormatFloat(instance.Stats.FreeSpace.Average, 'f', -1, 64),
		strconv.FormatFloat(instance.Stats.FreeSpace.Minimum, 'f', -1, 64),
		strconv.FormatFloat(instance.Stats.FreeSpace.Maximum, 'f', -1, 64),
		strconv.FormatFloat(instance.Stats.Cpu.Average, 'f', -1, 64),
		strconv.FormatFloat(instance.Stats.Cpu.Peak, 'f', -1, 64),
	}
}

func getRdsUsageReport(ctx context.Context, aa aws.AwsAccount, tx *sql.Tx) (data [][]string, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]string, 0)
	data = append(data, rdsInstanceFormat)

	date, _ := history.GetHistoryDate()

	identity, err := aa.GetAwsAccountIdentity()
	if err != nil {
		return
	}

	user, err := users.GetUserWithId(tx, aa.UserId)
	if err != nil {
		return
	}

	parameters := rds.RdsQueryParams{
		AccountList: []string{identity},
		Date: date,
	}

	logger.Debug("Getting RDS Usage Report for account", map[string]interface{}{
		"account": aa,
	})
	_, reports, err := rds.GetRdsData(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatRdsInstance(report.Instance)
			data = append(data, row)
		}
	}
	return
}
