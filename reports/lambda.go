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
	"github.com/trackit/trackit-server/usageReports/lambda"
	"github.com/trackit/trackit-server/users"
)

var lambdaFunctionFormat = [][]cell{{
	newCell("", 6).addStyle(textCenter, backgroundGrey),
	newCell("Invocations", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Duration", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("", 1).addStyle(textCenter, backgroundGrey),
}, {
	newCell("Account").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Name").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Version").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Runtime").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Size (Bytes)").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Memory (MegaBytes)").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Total").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Failed").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Average").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Maximum").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Tags").addStyle(textCenter, textBold, backgroundGrey),
}}

func formatLambdaFunction(report lambda.FunctionReport) []cell {
	function := report.Function
	tags := formatTags(function.Tags)
	return []cell{
		newCell(report.Account),
		newCell(function.Name),
		newCell(function.Version),
		newCell(function.Runtime),
		newCell(function.Size),
		newCell(function.Memory),
		newCell(formatMetric(function.Stats.Invocations.Total)),
		newCell(formatMetric(function.Stats.Invocations.Failed)),
		newCell(formatMetric(function.Stats.Duration.Average)),
		newCell(formatMetric(function.Stats.Duration.Maximum)),
		newCell(strings.Join(tags, ";")),
	}
}

func getLambdaUsageReport(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0, len(lambdaFunctionFormat))
	for _, headerRow := range lambdaFunctionFormat {
		data = append(data, headerRow)
	}

	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}

	if len(aas) < 1 {
		err = errors.New("Missing AWS Account for Lambda Usage Report")
		return
	}

	identities := getIdentities(aas)

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
	})
	_, reports, err := lambda.GetLambdaData(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatLambdaFunction(report)
			data = append(data, row)
		}
	}
	return
}
