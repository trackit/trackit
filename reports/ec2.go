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
	"strconv"
	"strings"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/ec2"
	"github.com/trackit/trackit-server/users"
)

var ec2InstanceFormat = []cell{
	newCell("ID").addStyle(textBold, backgroundGreen),
	newCell("Name").addStyle(textBold),
	newCell("Type").addStyle(textBold),
	newCell("Region").addStyle(textBold),
	newCell("Purchasing option").addStyle(textBold),
	newCell("Cost", 1).addStyle(textBold),
	newCell("CPU (Percentage)", 2).addStyle(textBold),
	newCell("Network (Bytes)", 2).addStyle(textBold),
	newCell("I/O (Bytes)", 4).addStyle(textBold),
	newCell("Key Pair").addStyle(textBold),
	newCell("Tags").addStyle(textBold),
	newCell("Name").addStyle(textBold),
	//	{"CPU Average (Percentage)", 1},
	//	{"CPU Peak (Percentage)", 1},
	//	{"Network In (Bytes)", 1},
	//	{"Network Out (Bytes)", 1},
	//	{"I/O Read (Bytes)", 1},
	//	{"I/O Read - Detailed (Bytes)", 1},
	//	{"I/O Write (Bytes)", 1},
	//	{"I/O Write - Detailed (Bytes)", 1},
}

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
	ioReadDetails := make([]string, 0)
	ioWriteDetails := make([]string, 0)
	for volume, size := range instance.Stats.Volumes.Read {
		ioRead += int(size)
		ioReadDetails = append(ioReadDetails, fmt.Sprintf("%s:%d", volume, int(size)))
	}
	for volume, size := range instance.Stats.Volumes.Write {
		ioWrite += int(size)
		ioWriteDetails = append(ioWriteDetails, fmt.Sprintf("%s:%d", volume, int(size)))
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
		newCell(strconv.FormatFloat(cost, 'f', -1, 64)),
		newCell(strconv.FormatFloat(instance.Stats.Cpu.Average, 'f', -1, 64)),
		newCell(strconv.FormatFloat(instance.Stats.Cpu.Peak, 'f', -1, 64)),
		newCell(strconv.FormatFloat(instance.Stats.Network.In, 'f', -1, 64)),
		newCell(strconv.FormatFloat(instance.Stats.Network.Out, 'f', -1, 64)),
		newCell(strconv.Itoa(ioRead)),
		newCell(strings.Join(ioReadDetails, ";")),
		newCell(strconv.Itoa(ioWrite)),
		newCell(strings.Join(ioWriteDetails, ";")),
		newCell(instance.KeyPair),
		newCell(strings.Join(tags, ";")),
	}
}

func getEc2UsageReport(ctx context.Context, aa aws.AwsAccount, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0)
	data = append(data, ec2InstanceFormat)

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
