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
	"errors"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/rds"
	"github.com/trackit/trackit-server/users"
)

var rdsInstanceFormat = [][]cell{{
	newCell("", 8).addStyle(textCenter, backgroundGrey),
	newCell("Storage - Available (Bytes)", 3).addStyle(textCenter, textBold, backgroundGrey),
	newCell("CPU (Percentage)", 2).addStyle(textCenter, textBold, backgroundGrey),
}, {
	newCell("Account").addStyle(textCenter, textBold, backgroundGrey),
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

func formatRdsInstance(report rds.InstanceReport) []cell {
	instance := report.Instance
	return []cell{
		newCell(report.Account),
		newCell(instance.DBInstanceIdentifier),
		newCell(instance.DBInstanceClass),
		newCell(instance.AvailabilityZone),
		newCell(getTotal(instance.Costs)),
		newCell(instance.Engine),
		newCell(instance.MultiAZ),
		newCell(instance.AllocatedStorage),
		newCell(formatMetric(instance.Stats.FreeSpace.Average)),
		newCell(formatMetric(instance.Stats.FreeSpace.Minimum)),
		newCell(formatMetric(instance.Stats.FreeSpace.Maximum)),
		newCell(formatMetricPercentage(instance.Stats.Cpu.Average)),
		newCell(formatMetricPercentage(instance.Stats.Cpu.Peak)),
	}
}

func getRdsUsageReport(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0, len(rdsInstanceFormat))
	for _, headerRow := range rdsInstanceFormat {
		data = append(data, headerRow)
	}

	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}

	if len(aas) < 1 {
		err = errors.New("Missing AWS Account for RDS Usage Report")
		return
	}

	identities := getIdentities(aas)

	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}

	parameters := rds.RdsQueryParams{
		AccountList: identities,
		Date:        date,
	}

	logger.Debug("Getting RDS Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
	})
	_, reports, err := rds.GetRdsData(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatRdsInstance(report)
			data = append(data, row)
		}
	}
	return
}
