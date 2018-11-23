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
	"github.com/trackit/trackit-server/usageReports/ec2"
	"github.com/trackit/trackit-server/users"
)

var ec2InstanceFormat = [][]cell{{
	newCell("", 6).addStyle(textCenter, backgroundGrey),
	newCell("CPU (Percentage)", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Network (Bytes)", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("I/O (Bytes)", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("", 2).addStyle(textCenter, backgroundGrey),
}, {
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

func formatEc2Instance(instance ec2.Instance) []cell {
	var cost float64
	for _, value := range instance.Costs {
		cost += value
	}
	name := ""
	if value, ok := instance.Tags["Name"]; ok {
		name = value
	}
	ioRead, ioWrite := 0, 0
	for _, size := range instance.Stats.Volumes.Read {
		ioRead += int(size)
	}
	for _, size := range instance.Stats.Volumes.Write {
		ioWrite += int(size)
	}
	tags := make([]string, 0)
	for key, value := range instance.Tags {
		tags = append(tags, fmt.Sprintf("%s:%s", key, value))
	}
	return []cell{
		newCell(instance.Id),
		newCell(name),
		newCell(instance.Type),
		newCell(instance.Region),
		newCell(instance.Purchasing),
		newCell(cost),
		newCell(instance.Stats.Cpu.Average / 100),
		newCell(instance.Stats.Cpu.Peak / 100),
		newCell(instance.Stats.Network.In),
		newCell(instance.Stats.Network.Out),
		newCell(ioRead),
		newCell(ioWrite),
		newCell(instance.KeyPair),
		newCell(strings.Join(tags, ";")),
	}
}

func getEc2UsageReport(ctx context.Context, aa aws.AwsAccount, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0)
	for _, headerRow := range ec2InstanceFormat {
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

	parameters := ec2.Ec2QueryParams{
		AccountList: []string{identity},
		Date:        date,
	}

	logger.Debug("Getting EC2 Usage Report for account", map[string]interface{}{
		"account": aa,
	})
	_, reports, err := ec2.GetEc2Data(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatEc2Instance(report.Instance)
			data = append(data, row)
		}
	}
	return
}
