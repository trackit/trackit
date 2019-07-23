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
	"errors"
	"strings"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/ec2"
	"github.com/trackit/trackit-server/users"
)

var ec2InstanceFormat = [][]cell{{
	newCell("", 7).addStyle(textCenter, backgroundGrey),
	newCell("CPU (Percentage)", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Network (Bytes)", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("I/O (Bytes)", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("", 2).addStyle(textCenter, backgroundGrey),
}, {
	newCell("Account").addStyle(textCenter, textBold, backgroundGrey),
	newCell("ID").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Name").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Type").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Region").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Purchasing option").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Cost").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Average").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Peak").addStyle(textCenter, textBold, backgroundGrey),
	newCell("In").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Out").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Read").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Write").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Key Pair").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Tags").addStyle(textCenter, textBold, backgroundGrey),
}}

func formatEc2Instance(report ec2.InstanceReport) []cell {
	instance := report.Instance
	name := ""
	if value, ok := instance.Tags["Name"]; ok {
		name = value
	}
	tags := formatTags(instance.Tags)
	return []cell{
		newCell(report.Account),
		newCell(instance.Id),
		newCell(name),
		newCell(instance.Type),
		newCell(instance.Region),
		newCell(instance.Purchasing),
		newCell(getTotal(instance.Costs)),
		newCell(formatMetricPercentage(instance.Stats.Cpu.Average)),
		newCell(formatMetricPercentage(instance.Stats.Cpu.Peak)),
		newCell(formatMetric(instance.Stats.Network.In)),
		newCell(formatMetric(instance.Stats.Network.Out)),
		newCell(getTotal(instance.Stats.Volumes.Read)),
		newCell(getTotal(instance.Stats.Volumes.Write)),
		newCell(instance.KeyPair),
		newCell(strings.Join(tags, ";")),
	}
}

func getEc2UsageReport(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0, len(ec2InstanceFormat))
	for _, headerRow := range ec2InstanceFormat {
		data = append(data, headerRow)
	}

	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}

	if len(aas) < 1 {
		err = errors.New("Missing AWS Account for EC2 Usage Report")
		return
	}

	identities := getIdentities(aas)

	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}

	parameters := ec2.Ec2QueryParams{
		AccountList: identities,
		Date:        date,
	}

	logger.Debug("Getting EC2 Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
	})
	_, reports, err := ec2.GetEc2Data(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatEc2Instance(report)
			data = append(data, row)
		}
	}
	return
}
