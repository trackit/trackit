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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/usageReports/lambda"
	"github.com/trackit/trackit/users"
)

const lambdaUsageReportSheetName = "Lambda Usage Report"

var lambdaUsageReportModule = module{
	Name:          "Lambda Usage Report",
	SheetName:     lambdaUsageReportSheetName,
	ErrorName:     "lambdaUsageReportError",
	GenerateSheet: generateLambdaUsageReportSheet,
}

// generateLambdaUsageReportSheet will generate a sheet with Lambda usage report
// It will get data for given AWS account and for a given date
func generateLambdaUsageReportSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return lambdaUsageReportGenerateSheet(ctx, aas, date, tx, file)
}

func lambdaUsageReportGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	data, err := lambdaUsageReportGetData(ctx, aas, date, tx)
	if err == nil {
		return lambdaUsageReportInsertDataInSheet(aas, file, data)
	}
	return
}

func lambdaUsageReportGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (reports []lambda.FunctionReport, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identities := getAwsIdentities(aas)
	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return
	}
	parameters := lambda.LambdaQueryParams{
		AccountList: identities,
		Date:        date,
	}
	logger.Debug("Getting Lambda Usage Report for accounts", map[string]interface{}{
		"accounts": aas,
		"date":     date,
	})
	_, reports, err = lambda.GetLambdaData(ctx, parameters, user, tx)
	if err != nil {
		logger.Error("An error occurred while generating an Lambda Usage Report", map[string]interface{}{
			"error":    err,
			"accounts": aas,
			"date":     date,
		})
	}
	return
}

func lambdaUsageReportInsertDataInSheet(aas []aws.AwsAccount, file *excelize.File, data []lambda.FunctionReport) (err error) {
	file.NewSheet(lambdaUsageReportSheetName)
	lambdaUsageReportGenerateHeader(file)
	line := 3
	for _, report := range data {
		account := getAwsAccount(report.Account, aas)
		formattedAccount := report.Account
		if account != nil {
			formattedAccount = formatAwsAccount(*account)
		}
		function := report.Function
		tags := formatTags(function.Tags)
		totalCol := "G" + strconv.Itoa(line)
		failedCol := "H" + strconv.Itoa(line)
		formula := fmt.Sprintf(`IF(%s="N/A","",%s/%s)`, totalCol, failedCol, totalCol)
		invocationsFailedPercentage := newFormula(formula, "I"+strconv.Itoa(line)).addStyles("percentage")
		invocationsFailedPercentage = invocationsFailedPercentage.addConditionalFormat("validPercentage", "red", "borders")
		invocationsFailedPercentage = invocationsFailedPercentage.addConditionalFormat("above90percent", "green", "borders")
		invocationsFailedPercentage = invocationsFailedPercentage.addConditionalFormat("above75percent", "orange", "borders")
		cells := cells{
			newCell(formattedAccount, "A"+strconv.Itoa(line)),
			newCell(function.Name, "B"+strconv.Itoa(line)),
			newCell(function.Version, "C"+strconv.Itoa(line)),
			newCell(function.Runtime, "D"+strconv.Itoa(line)),
			newCell(function.Size, "E"+strconv.Itoa(line)),
			newCell(function.Memory, "F"+strconv.Itoa(line)),
			newCell(formatMetric(function.Stats.Invocations.Total), "G"+strconv.Itoa(line)),
			newCell(formatMetric(function.Stats.Invocations.Failed), "H"+strconv.Itoa(line)),
			invocationsFailedPercentage,
			newCell(formatMetric(function.Stats.Duration.Average), "J"+strconv.Itoa(line)),
			newCell(formatMetric(function.Stats.Duration.Maximum), "K"+strconv.Itoa(line)),
			newCell(strings.Join(tags, ";"), "L"+strconv.Itoa(line)),
		}
		cells.addStyles("borders", "centerText").setValues(file, lambdaUsageReportSheetName)
		line++
	}
	return
}

func lambdaUsageReportGenerateHeader(file *excelize.File) {
	header := cells{
		newCell("Account", "A1").mergeTo("A2"),
		newCell("Name", "B1").mergeTo("B2"),
		newCell("Version", "C1").mergeTo("C2"),
		newCell("Runtime", "D1").mergeTo("D2"),
		newCell("Size (Bytes)", "E1").mergeTo("E2"),
		newCell("Memory (MegaBytes)", "F1").mergeTo("F2"),
		newCell("Invocations", "G1").mergeTo("I2"),
		newCell("Total", "G2"),
		newCell("Failed", "H2"),
		newCell("Success (Percentage)", "I2"),
		newCell("Duration (Milliseconds)", "J1").mergeTo("K1"),
		newCell("Average", "J2"),
		newCell("Maximum", "K2"),
		newCell("Tags", "L1").mergeTo("L2"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, lambdaUsageReportSheetName)
	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("B", 40),
		newColumnWidth("C", 15).toColumn("I"),
		newColumnWidth("F", 17.5),
		newColumnWidth("J", 12.5).toColumn("J"),
		newColumnWidth("L", 30),
	}
	columns.setValues(file, lambdaUsageReportSheetName)
}
