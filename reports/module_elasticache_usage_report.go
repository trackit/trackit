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
	"github.com/trackit/trackit/usageReports/elasticache"
	"github.com/trackit/trackit/users"
)

const elastiCacheUsageReportSheetName = "ElastiCache Usage Report"

var elastiCacheUsageReportModule = module{
	Name:          "ElastiCache Usage Report",
	SheetName:     elastiCacheUsageReportSheetName,
	ErrorName:     "elastiCacheUsageReportError",
	GenerateSheet: generateElastiCacheUsageReportSheet,
}

// generateElastiCacheUsageReportSheet will generate a sheet with ElastiCache usage report
// It will get data for given AWS account and for a given date
func generateElastiCacheUsageReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return elastiCacheUsageReportGenerateSheet(ctx, aas, date, tx, file)
}

func elastiCacheUsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := elastiCacheUsageReportGetData(ctx, aas, date, tx)
	if err == nil {
		return elastiCacheUsageReportInsertDataInSheet(aas, file, data)
	}
	return
}

func elastiCacheUsageReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []elasticache.InstanceReport, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}
	parameters := elasticache.ElastiCacheQueryParams{
		AccountList: identities,
		Date:        date,
	}
	logger.Debug("Getting ElastiCache Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date,
	})
	_, reports, err = elasticache.GetElastiCacheData(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an ElastiCache Usage Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func elastiCacheUsageReportInsertDataInSheet(aas []aws.AwsAccount, file *excelize.File, data []elasticache.InstanceReport) (err error) {
	file.NewSheet(elastiCacheUsageReportSheetName)
	elastiCacheUsageReportGenerateHeader(file)
	line := 3
	for _, report := range data {
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		instance := report.Instance
		tags := formatTags(instance.Tags)
		cpuAverage := newCell(formatMetricPercentage(instance.Stats.Cpu.Average), "H"+strconv.Itoa(line)).addStyles("percentage")
		cpuAverage = cpuAverage.addConditionalFormat("validPercentage", "green", "borders")
		cpuAverage = cpuAverage.addConditionalFormat("above60percent", "red", "borders")
		cpuAverage = cpuAverage.addConditionalFormat("above30percent", "orange", "borders")
		cpuPeak := newCell(formatMetricPercentage(instance.Stats.Cpu.Peak), "I"+strconv.Itoa(line)).addStyles("percentage")
		cpuPeak = cpuPeak.addConditionalFormat("validPercentage", "green", "borders")
		cpuPeak = cpuPeak.addConditionalFormat("above80percent", "red", "borders")
		cpuPeak = cpuPeak.addConditionalFormat("above60percent", "orange", "borders")
		cells := cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)),
			newCell(instance.Id, "B"+strconv.Itoa(line)),
			newCell(instance.NodeType, "C"+strconv.Itoa(line)),
			newCell(instance.Region, "D"+strconv.Itoa(line)),
			newCell(getTotal(instance.Costs), "E"+strconv.Itoa(line)).addStyles("price"),
			newCell(instance.Engine, "F"+strconv.Itoa(line)),
			newCell(instance.EngineVersion, "G"+strconv.Itoa(line)),
			cpuAverage,
			cpuPeak,
			newCell(formatMetric(instance.Stats.Network.In), "J"+strconv.Itoa(line)),
			newCell(formatMetric(instance.Stats.Network.Out), "K"+strconv.Itoa(line)),
			newCell(strings.Join(tags, ";"), "L"+strconv.Itoa(line)),
		}
		cells.addStyles("borders", "centerText").setValues(file, elastiCacheUsageReportSheetName)
		line++
	}
	return
}

func elastiCacheUsageReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A2"),
		newCell("ID", "B1").mergeTo("B2"),
		newCell("Type", "C1").mergeTo("C2"),
		newCell("Region", "D1").mergeTo("D2"),
		newCell("Cost", "E1").mergeTo("E2"),
		newCell("Engine", "F1").mergeTo("G1"),
		newCell("Name", "F2"),
		newCell("Version", "G2"),
		newCell("CPU (Percentage)", "H1").mergeTo("I1"),
		newCell("Average", "H2"),
		newCell("Peak", "I2"),
		newCell("Network (Bytes)", "J1").mergeTo("K1"),
		newCell("In", "J2"),
		newCell("Out", "K2"),
		newCell("Tags", "L1").mergeTo("L2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, elastiCacheUsageReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 35).toColumn("C"),
		newColumnWidth("D", 15).toColumn("E"),
		newColumnWidth("F", 20).toColumn("G"),
		newColumnWidth("H", 12.5).toColumn("N"),
		newColumnWidth("L", 30),
	}
	columns.setValues(file, elastiCacheUsageReportSheetName)
}
