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
	"github.com/trackit/trackit-server/usageReports/lambda"
	"github.com/trackit/trackit-server/users"
)

var lambdaFunctionFormat = [][]cell{{
	newCell("", 5).addStyle(textCenter, backgroundGrey),
	newCell("Invocations", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("Duration", 2).addStyle(textCenter, textBold, backgroundGrey),
	newCell("", 1).addStyle(textCenter, backgroundGrey),
}, {
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

func formatLambdaFunction(function lambda.Function) []cell {
	tags := make([]string, 0)
	for key, value := range function.Tags {
		tags = append(tags, fmt.Sprintf("%s:%s", key, value))
	}
	return []cell{
		newCell(function.Name),
		newCell(function.Version),
		newCell(function.Runtime),
		newCell(function.Size),
		newCell(function.Memory),
		newCell(function.Stats.Invocations.Total),
		newCell(function.Stats.Invocations.Failed),
		newCell(function.Stats.Duration.Average),
		newCell(function.Stats.Duration.Maximum),
		newCell(strings.Join(tags, ";")),
	}
}

func getLambdaUsageReport(ctx context.Context, aa aws.AwsAccount, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	data = make([][]cell, 0)
	for _, headerRow := range lambdaFunctionFormat {
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

	parameters := lambda.LambdaQueryParams{
		AccountList: []string{identity},
		Date:        date,
	}

	logger.Debug("Getting Lambda Usage Report for account", map[string]interface{}{
		"account": aa,
	})
	_, reports, err := lambda.GetLambdaData(ctx, parameters, user, tx)
	if err != nil {
		return
	}

	if reports != nil && len(reports) > 0 {
		for _, report := range reports {
			row := formatLambdaFunction(report.Function)
			data = append(data, row)
		}
	}
	return
}
