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
	"strings"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/elasticache"
	"github.com/trackit/trackit-server/users"
)

var elasticacheInstanceFormat = [][]cell{{
	newCell("", 5).addStyle(textCenter, backgroundGrey),
	newCell("Engine", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("CPU (Percentage)", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Network (Bytes)", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("", 1).addStyle(textCenter, backgroundGrey),
}, {
	newCell("Account").addStyle(textCenter, textBold, backgroundGrey),
	newCell("ID").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Type").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Region").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Cost").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Name").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Version").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Average").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Peak").addStyle(textCenter, textBold, backgroundGrey),
	newCell("In").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Out").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Tags").addStyle(textCenter, textBold, backgroundGrey),
}}

func formatElasticacheInstance(report elasticache.InstanceReport) []cell {
	instance := report.Instance
	tags := formatTags(instance.Tags)
	return []cell{
		newCell(report.Account),
		newCell(instance.Id),
		newCell(instance.NodeType),
		newCell(instance.NodeType),
		newCell(instance.Region),
		newCell(getTotal(instance.Costs)),
		newCell(instance.Engine),
		newCell(instance.EngineVersion),
		newCell(formatMetricPercentage(instance.Stats.Cpu.Average)),
		newCell(formatMetricPercentage(instance.Stats.Cpu.Peak)),
		newCell(formatMetric(instance.Stats.Network.In)),
		newCell(formatMetric(instance.Stats.Network.Out)),
		newCell(strings.Join(tags, ";")),
	}
}

func getElasticacheUsageReport(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0)
	for _, headerRow := range elasticacheInstanceFormat {
		data = append(data, headerRow)
	}

	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}

	if len(aas) < 1 {
		err = errors.New("Missing AWS Account for Elasticache Usage Report")
		return
	}

	identities := make([]string, 0)
	for _, account := range aas {
		identities = append(identities, account.AwsIdentity)
	}

	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}

	parameters := elasticache.ElastiCacheQueryParams{
		AccountList: identities,
		Date:        date,
	}

	logger.Debug("Getting Elasticache Usage Report for account", map[string]interface{}{
		"accounts": aas,
	})
	_, reports, err := elasticache.GetElastiCacheData(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatElasticacheInstance(report)
			data = append(data, row)
		}
	}
	return
}
