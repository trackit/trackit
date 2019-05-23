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
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/rds"
	"github.com/trackit/trackit-server/users"
)

const rdsUsageReportSheetName = "RDS Usage Report"

var rdsUsageReportModule = module{
	Name:          "RDS Usage Report",
	SheetName:     rdsUsageReportSheetName,
	ErrorName:     "rdsUsageReportError",
	GenerateSheet: generateRdsUsageReportSheet,
}

func generateRdsUsageReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return rdsUsageReportGenerateSheet(ctx, aas, date, tx, file)
}

func rdsUsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := rdsUsageReportGetData(ctx, aas, date, tx)
	if err == nil {
		return rdsUsageReportInsertDataInSheet(ctx, aas, file, data)
	} else {
		return
	}
}

func rdsUsageReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []rds.InstanceReport, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	identities := getAwsIdentities(aas)

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
		"date":     date,
	})
	_, reports, err = rds.GetRdsData(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an RDS Usage Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func rdsUsageReportInsertDataInSheet(_ context.Context, aas []aws.AwsAccount, file *excelize.File, data []rds.InstanceReport) (err error) {
	file.NewSheet(rdsUsageReportSheetName)
	rdsUsageReportGenerateHeader(file)
	line := 4
	for _, report := range data {
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		instance := report.Instance
		cpuAverage := newCell(formatMetricPercentage(instance.Stats.Cpu.Average), "L"+strconv.Itoa(line)).addStyles("percentage")
		cpuAverage = cpuAverage.addConditionalFormat("validPercentage", "green", "borders")
		cpuAverage = cpuAverage.addConditionalFormat("above60percent", "red", "borders")
		cpuAverage = cpuAverage.addConditionalFormat("above30percent", "orange", "borders")
		cpuPeak := newCell(formatMetricPercentage(instance.Stats.Cpu.Peak), "M"+strconv.Itoa(line)).addStyles("percentage")
		cpuPeak = cpuPeak.addConditionalFormat("validPercentage", "green", "borders")
		cpuPeak = cpuPeak.addConditionalFormat("above80percent", "red", "borders")
		cpuPeak = cpuPeak.addConditionalFormat("above60percent", "orange", "borders")
		cells := cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)),
			newCell(instance.DBInstanceIdentifier, "B"+strconv.Itoa(line)),
			newCell(instance.DBInstanceClass, "C"+strconv.Itoa(line)),
			newCell(instance.AvailabilityZone, "D"+strconv.Itoa(line)),
			newCell(getTotal(instance.Costs), "E"+strconv.Itoa(line)).addStyles("price"),
			newCell(instance.Engine, "F"+strconv.Itoa(line)),
			newCell(instance.MultiAZ, "G"+strconv.Itoa(line)),
			newCell(instance.AllocatedStorage, "H"+strconv.Itoa(line)),
			newCell(formatMetric(instance.Stats.FreeSpace.Average), "I"+strconv.Itoa(line)),
			newCell(formatMetric(instance.Stats.FreeSpace.Minimum), "J"+strconv.Itoa(line)),
			newCell(formatMetric(instance.Stats.FreeSpace.Maximum), "K"+strconv.Itoa(line)),
			cpuAverage,
			cpuPeak,
		}
		cells.addStyles("borders", "centerText").setValues(file, rdsUsageReportSheetName)
		line++
	}
	return
}

func rdsUsageReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A3"),
		newCell("Name", "B1").mergeTo("B3"),
		newCell("Type", "C1").mergeTo("C3"),
		newCell("Region", "D1").mergeTo("D3"),
		newCell("Cost", "E1").mergeTo("E3"),
		newCell("Engine", "F1").mergeTo("F3"),
		newCell("Multi A-Z", "G1").mergeTo("G3"),
		newCell("Storage", "H1").mergeTo("K1"),
		newCell("Total (GigaBytes)", "H2").mergeTo("H3"),
		newCell("Available (Bytes)", "I2").mergeTo("K2"),
		newCell("Average", "I3"),
		newCell("Minimum", "J3"),
		newCell("Maximum", "K3"),
		newCell("CPU (Percentage)", "L1").mergeTo("M1"),
		newCell("Average", "L2").mergeTo("L3"),
		newCell("Peak", "M2").mergeTo("M3"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, rdsUsageReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 35),
		newColumnWidth("C", 15).toColumn("K"),
		newColumnWidth("L", 12.5).toColumn("M"),
	}
	columns.setValues(file, rdsUsageReportSheetName)
	return
}
