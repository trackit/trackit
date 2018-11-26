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

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/rds"
	"github.com/trackit/trackit-server/users"
)

var rdsInstanceFormat = [][]cell{{
	newCell("", 7).addStyle(textCenter, backgroundGrey),
	newCell("Storage - Available (Bytes)", 3).addStyle(textCenter, textBold, backgroundGrey),
	newCell("CPU (Percentage)", 2).addStyle(textCenter, textBold, backgroundGrey),
}, {
	newCell("Name").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Type").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Region").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Cost").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Engine").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Multi A-Z").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Storage - Total (GigaBytes)").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Average").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Minimum").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Maximum").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Average").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Peak").addStyle(textCenter, textBold, backgroundGrey),
}}

func formatRdsInstance(instance rds.Instance) []cell {
	var cost float64

	for _, value := range instance.Costs {
		cost += value
	}

	return []cell{
		newCell(instance.DBInstanceIdentifier),
		newCell(instance.DBInstanceClass),
		newCell(instance.AvailabilityZone),
		newCell(cost),
		newCell(instance.Engine),
		newCell(instance.MultiAZ),
		newCell(instance.AllocatedStorage),
		newCell(instance.Stats.FreeSpace.Average),
		newCell(instance.Stats.FreeSpace.Minimum),
		newCell(instance.Stats.FreeSpace.Maximum),
		newCell(instance.Stats.Cpu.Average),
		newCell(instance.Stats.Cpu.Peak),
	}
}

func getRdsUsageReport(ctx context.Context, aa aws.AwsAccount, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0)
	for _, headerRow := range rdsInstanceFormat {
		data = append(data, headerRow)
	}

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
		Date:        date,
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
