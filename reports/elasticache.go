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
	"fmt"
	"strings"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/elasticache"
	"github.com/trackit/trackit-server/users"
)

var elasticacheInstanceFormat = [][]cell{{
	newCell("", 6).addStyle(textCenter, backgroundGrey),
	newCell("Storage", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("CPU (Percentage)", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Memory Pressure (Percentage)", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("", 1).addStyle(textCenter, backgroundGrey),
}, {
	newCell("ID").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Name").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Type").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Region").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Instances").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Cost").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total (GigaBytes)").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Unused (MegaBytes)").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Average").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Peak").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Average").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Peak").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Tags").addStyle(textCenter, textBold, backgroundGrey),
}}

func formatElasticacheInstance(instance elasticache.Instance) []cell {
	var cost float64
	for _, value := range instance.Costs {
		cost += value
	}
	tags := make([]string, 0)
	for key, value := range instance.Tags {
		tags = append(tags, fmt.Sprintf("%s:%s", key, value))
	}
	return []cell{
/*		newCell(instance.DomainID),
		newCell(instance.DomainName),
		newCell(instance.InstanceType),*/
		newCell(instance.Region),
//		newCell(instance.InstanceCount),
		newCell(cost),
//		newCell(instance.TotalStorageSpace),
//		newCell(instance.Stats.FreeSpace),
		newCell(instance.Stats.Cpu.Average),
		newCell(instance.Stats.Cpu.Peak),
//		newCell(instance.Stats.JVMMemoryPressure.Average),
//		newCell(instance.Stats.JVMMemoryPressure.Peak),
		newCell(strings.Join(tags, ";")),
	}
}

func getElasticacheUsageReport(ctx context.Context, aa aws.AwsAccount, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0)
	for _, headerRow := range elasticacheInstanceFormat {
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

	parameters := elasticache.ElastiCacheQueryParams{
		AccountList: []string{identity},
		Date:        date,
	}

	logger.Debug("Getting Elasticache Usage Report for account", map[string]interface{}{
		"account": aa,
	})
	_, reports, err := elasticache.GetElastiCacheData(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatElasticacheInstance(report.Instance)
			data = append(data, row)
		}
	}
	return
}
