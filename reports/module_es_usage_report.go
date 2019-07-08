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

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/usageReports/es"
	"github.com/trackit/trackit-server/users"
)

const esUsageReportSheetName = "ES Usage Report"

var esUsageReportModule = module{
	Name:          "ES Usage Report",
	SheetName:     esUsageReportSheetName,
	ErrorName:     "esUsageReportError",
	GenerateSheet: generateEsUsageReportSheet,
}

// generateEsUsageReportSheet will generate a sheet with ES usage report
// It will get data for given AWS account and for a given date
func generateEsUsageReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return esUsageReportGenerateSheet(ctx, aas, date, tx, file)
}

func esUsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := esUsageReportGetData(ctx, aas, date, tx)
	if err == nil {
		return esUsageReportInsertDataInSheet(aas, file, data)
	}
	return
}

func esUsageReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []es.DomainReport, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}
	parameters := es.EsQueryParams{
		AccountList: identities,
		Date:        date,
	}
	logger.Debug("Getting ES Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date,
	})
	_, reports, err = es.GetEsData(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an ES Usage Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func esUsageReportInsertDataInSheet(aas []aws.AwsAccount, file *excelize.File, data []es.DomainReport) (err error) {
	file.NewSheet(esUsageReportSheetName)
	esUsageReportGenerateHeader(file)
	line := 3
	for _, report := range data {
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		domain := report.Domain
		tags := formatTags(domain.Tags)
		cpuAverage := newCell(formatMetricPercentage(domain.Stats.Cpu.Average), "J"+strconv.Itoa(line)).addStyles("percentage")
		cpuAverage = cpuAverage.addConditionalFormat("validPercentage", "green", "borders")
		cpuAverage = cpuAverage.addConditionalFormat("above60percent", "red", "borders")
		cpuAverage = cpuAverage.addConditionalFormat("above30percent", "orange", "borders")
		cpuPeak := newCell(formatMetricPercentage(domain.Stats.Cpu.Peak), "K"+strconv.Itoa(line)).addStyles("percentage")
		cpuPeak = cpuPeak.addConditionalFormat("validPercentage", "green", "borders")
		cpuPeak = cpuPeak.addConditionalFormat("above80percent", "red", "borders")
		cpuPeak = cpuPeak.addConditionalFormat("above60percent", "orange", "borders")
		memoryAverage := newCell(formatMetricPercentage(domain.Stats.JVMMemoryPressure.Average), "L"+strconv.Itoa(line)).addStyles("percentage")
		memoryAverage = memoryAverage.addConditionalFormat("validPercentage", "green", "borders")
		memoryAverage = memoryAverage.addConditionalFormat("above85percent", "red", "borders")
		memoryAverage = memoryAverage.addConditionalFormat("above75percent", "orange", "borders")
		memoryPeak := newCell(formatMetricPercentage(domain.Stats.JVMMemoryPressure.Peak), "M"+strconv.Itoa(line)).addStyles("percentage")
		memoryPeak = memoryPeak.addConditionalFormat("validPercentage", "green", "borders")
		memoryPeak = memoryPeak.addConditionalFormat("above85percent", "red", "borders")
		memoryPeak = memoryPeak.addConditionalFormat("above75percent", "orange", "borders")
		cells := cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)),
			newCell(domain.DomainID, "B"+strconv.Itoa(line)),
			newCell(domain.DomainName, "C"+strconv.Itoa(line)),
			newCell(domain.InstanceType, "D"+strconv.Itoa(line)),
			newCell(domain.Region, "E"+strconv.Itoa(line)),
			newCell(domain.InstanceCount, "F"+strconv.Itoa(line)),
			newCell(getTotal(domain.Costs), "G"+strconv.Itoa(line)).addStyles("price"),
			newCell(domain.TotalStorageSpace, "H"+strconv.Itoa(line)),
			newCell(formatMetric(domain.Stats.FreeSpace), "I"+strconv.Itoa(line)),
			cpuAverage,
			cpuPeak,
			memoryAverage,
			memoryPeak,
			newCell(strings.Join(tags, ";"), "N"+strconv.Itoa(line)),
		}
		cells.addStyles("borders", "centerText").setValues(file, esUsageReportSheetName)
		line++
	}
	return
}

func esUsageReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A2"),
		newCell("ID", "B1").mergeTo("B2"),
		newCell("Name", "C1").mergeTo("C2"),
		newCell("Type", "D1").mergeTo("D2"),
		newCell("Region", "E1").mergeTo("E2"),
		newCell("Instances", "F1").mergeTo("F2"),
		newCell("Cost", "G1").mergeTo("G2"),
		newCell("Storage", "H1").mergeTo("I1"),
		newCell("Total (GigaBytes)", "H2"),
		newCell("Unused (MegaBytes)", "I2"),
		newCell("CPU (Percentage)", "J1").mergeTo("K1"),
		newCell("Average", "J2"),
		newCell("Peak", "K2"),
		newCell("Memory Pressure (Percentage)", "L1").mergeTo("M1"),
		newCell("Average", "L2"),
		newCell("Peak", "M2"),
		newCell("Tags", "N1").mergeTo("N2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, esUsageReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 35).toColumn("C"),
		newColumnWidth("D", 20),
		newColumnWidth("E", 15),
		newColumnWidth("F", 12.5).toColumn("M"),
		newColumnWidth("H", 17.5).toColumn("I"),
		newColumnWidth("N", 30),
	}
	columns.setValues(file, esUsageReportSheetName)
	return
}
