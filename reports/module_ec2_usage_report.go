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
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/usageReports/ec2"
	"github.com/trackit/trackit/users"
)

const ec2UsageReportSheetName = "EC2 Usage Report"
const ec2SizingRecommendationsSheetName = "EC2 Sizing Recommendations"

var ec2UsageReportModule = module{
	Name:          "EC2 Usage Report",
	SheetName:     ec2UsageReportSheetName,
	ErrorName:     "ec2UsageReportError",
	GenerateSheet: generateEc2UsageReportSheet,
}

// generateEc2UsageReportSheet will generate a sheet with EC2 usage report
// It will get data for given AWS account and for a given date
func generateEc2UsageReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return ec2UsageReportGenerateSheet(ctx, aas, date, tx, file)
}

func ec2UsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := ec2UsageReportGetData(ctx, aas, date, tx)
	if err == nil {
		return ec2UsageReportInsertDataInSheet(aas, file, data)
	}
	return
}

func ec2UsageReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []ec2.InstanceReport, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identities := getAwsIdentities(aas)
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
		"date":     date,
	})
	_, reports, err = ec2.GetEc2Data(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an EC2 Usage Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func ec2UsageReportInsertDataInSheet(aas []aws.AwsAccount, file *excelize.File, data []ec2.InstanceReport) (err error) {
	file.NewSheet(ec2UsageReportSheetName)
	ec2SizingRecommendationsHeader(file)
	ec2UsageReportGenerateHeader(file)
	line := 3
	for _, report := range data {
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		instance := report.Instance
		name := ""
		if value, ok := instance.Tags["Name"]; ok {
			name = value
		}
		tags := formatTags(instance.Tags)
		cpuAverage := newCell(formatMetricPercentage(instance.Stats.Cpu.Average), "H"+strconv.Itoa(line)).addStyles("percentage")
		cpuAverage = cpuAverage.addConditionalFormat("validPercentage", "green", "borders")
		cpuAverage = cpuAverage.addConditionalFormat("above60percent", "red", "borders")
		cpuAverage = cpuAverage.addConditionalFormat("above30percent", "orange", "borders")
		cpuPeak := newCell(formatMetricPercentage(instance.Stats.Cpu.Peak), "I"+strconv.Itoa(line)).addStyles("percentage")
		cpuPeak = cpuPeak.addConditionalFormat("validPercentage", "green", "borders")
		cpuPeak = cpuPeak.addConditionalFormat("above80percent", "red", "borders")
		cpuPeak = cpuPeak.addConditionalFormat("above60percent", "orange", "borders")
		cellsReport := cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)),
			newCell(instance.Id, "B"+strconv.Itoa(line)),
			newCell(name, "C"+strconv.Itoa(line)),
			newCell(instance.Type, "D"+strconv.Itoa(line)),
			newCell(instance.Region, "E"+strconv.Itoa(line)),
			newCell(instance.Purchasing, "F"+strconv.Itoa(line)),
			newCell(getTotal(instance.Costs), "G"+strconv.Itoa(line)).addStyles("price"),
			cpuAverage,
			cpuPeak,
			newCell(formatMetric(instance.Stats.Network.In), "J"+strconv.Itoa(line)),
			newCell(formatMetric(instance.Stats.Network.Out), "K"+strconv.Itoa(line)),
			newCell(getTotal(instance.Stats.Volumes.Read), "L"+strconv.Itoa(line)),
			newCell(getTotal(instance.Stats.Volumes.Write), "M"+strconv.Itoa(line)),
			newCell(instance.KeyPair, "N"+strconv.Itoa(line)),
			newCell(strings.Join(tags, ";"), "O"+strconv.Itoa(line)),
		}
		cellsReport.addStyles("borders", "centerText").setValues(file, ec2UsageReportSheetName)
		ec2SizingRecommendationInsertData(file, instance, name, formattedAccount, line)
		line++
	}
	return
}

func ec2UsageReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A2"),
		newCell("ID", "B1").mergeTo("B2"),
		newCell("Name", "C1").mergeTo("C2"),
		newCell("Type", "D1").mergeTo("D2"),
		newCell("Region", "E1").mergeTo("E2"),
		newCell("Purchasing option", "F1").mergeTo("F2"),
		newCell("Cost", "G1").mergeTo("G2"),
		newCell("CPU (Percentage)", "H1").mergeTo("I1"),
		newCell("Average", "H2"),
		newCell("Peak", "I2"),
		newCell("Network (Bytes)", "J1").mergeTo("K1"),
		newCell("In", "J2"),
		newCell("Out", "K2"),
		newCell("I/O (Bytes)", "L1").mergeTo("M1"),
		newCell("Read", "L2"),
		newCell("Write", "M2"),
		newCell("Key Pair", "N1").mergeTo("N2"),
		newCell("Tags", "O1").mergeTo("O2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, ec2UsageReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 35).toColumn("C"),
		newColumnWidth("D", 15).toColumn("E"),
		newColumnWidth("F", 20),
		newColumnWidth("G", 12.5).toColumn("N"),
		newColumnWidth("L", 20).toColumn("M"),
		newColumnWidth("O", 30),
	}
	columns.setValues(file, ec2UsageReportSheetName)
	return
}

func ec2SizingRecommendationInsertData(file *excelize.File, instance ec2.Instance, name, formattedAccount string, line int)  {
	cellsRecommendation := cells{
		newCell(formattedAccount, "A"+strconv.Itoa(line)),
		newCell(instance.Id, "B"+strconv.Itoa(line)),
		newCell(name, "C"+strconv.Itoa(line)),
		newCell(instance.Region, "D"+strconv.Itoa(line)),
		newCell(instance.Type, "E"+strconv.Itoa(line)),
		newCell(instance.Recommendation.InstanceType, "F"+strconv.Itoa(line)),
		newCell(instance.Recommendation.Reason, "G"+strconv.Itoa(line)),
		newCell(instance.Recommendation.NewGeneration, "H"+strconv.Itoa(line)),
	}
	cellsRecommendation.addStyles("borders", "centerText").setValues(file, ec2SizingRecommendationsSheetName)
}

func ec2SizingRecommendationsHeader(file *excelize.File) {
	file.NewSheet(ec2SizingRecommendationsSheetName)
	header := cells{
		newCell("Account", "A1").mergeTo("A2"),
		newCell("ID", "B1").mergeTo("B2"),
		newCell("Name", "C1").mergeTo("C2"),
		newCell("Region", "D1").mergeTo("D2"),
		newCell("Type", "E1").mergeTo("E2"),
		newCell("Recommendation", "F1").mergeTo("H1"),
		newCell("Type", "F2"),
		newCell("Reason", "G2"),
		newCell("New Generations", "H2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, ec2SizingRecommendationsSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 35).toColumn("C"),
		newColumnWidth("D", 15).toColumn("E"),
		newColumnWidth("F", 20).toColumn("H"),
	}
	columns.setValues(file, ec2SizingRecommendationsSheetName)
	return
}
