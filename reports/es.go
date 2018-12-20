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
	"strings"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/es"
	"github.com/trackit/trackit-server/users"
)

var esDomainFormat = [][]cell{{
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

func formatEsDomain(domain es.Domain) []cell {
	tags := formatTags(domain.Tags)
	return []cell{
		newCell(domain.DomainID),
		newCell(domain.DomainName),
		newCell(domain.InstanceType),
		newCell(domain.Region),
		newCell(domain.InstanceCount),
		newCell(getTotal(domain.Costs)),
		newCell(domain.TotalStorageSpace),
		newCell(formatMetric(domain.Stats.FreeSpace)),
		newCell(formatMetricPercentage(domain.Stats.Cpu.Average)),
		newCell(formatMetricPercentage(domain.Stats.Cpu.Peak)),
		newCell(formatMetric(domain.Stats.JVMMemoryPressure.Average)),
		newCell(formatMetric(domain.Stats.JVMMemoryPressure.Peak)),
		newCell(strings.Join(tags, ";")),
	}
}

func getEsUsageReport(ctx context.Context, aa aws.AwsAccount, date time.Time, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0)
	for _, headerRow := range esDomainFormat {
		data = append(data, headerRow)
	}

	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}

	identity, err := aa.GetAwsAccountIdentity()
	if err != nil {
		return
	}

	user, err := users.GetUserWithId(tx, aa.UserId)
	if err != nil {
		return
	}

	parameters := es.EsQueryParams{
		AccountList: []string{identity},
		Date:        date,
	}

	logger.Debug("Getting ES Usage Report for account", map[string]interface{}{
		"account": aa,
	})
	_, reports, err := es.GetEsData(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatEsDomain(report.Domain)
			data = append(data, row)
		}
	}
	return
}
