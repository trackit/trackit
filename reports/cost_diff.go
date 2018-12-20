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
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports/history"
	"github.com/trackit/trackit-server/costs/diff"
)

type costDiffProduct map[time.Time]diff.PricePoint

var costDiffHeader = [][]cell{{
	newCell("Account").addStyle(textCenter, textBold, backgroundGrey),
	newCell("Product").addStyle(textCenter, textBold, backgroundGrey),
}}

func formatCostDiff(data []diff.PricePoint) (values costDiffProduct, err error) {
	values = make(costDiffProduct)
	for _, value := range data {
		date, err := time.Parse("2006-01-02T15:04:05.000Z", value.Date)
		if err != nil {
			return values, err
		}
		values[date] = value
	}
	return
}

func getValueForDate(values costDiffProduct, date time.Time) *diff.PricePoint {
	if value, ok := values[date] ; ok {
		return &value
	}
	return nil
}

func getCostDiff(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Getting Cost Differentiator Report for accounts", map[string]interface{}{
		"accounts": aas,
	})

	data = make([][]cell, 0)
	header := make([]cell, 0)
	header = append(header, costDiffHeader[0]...)

	var dateBegin, dateEnd time.Time
	if date.IsZero() {
		dateBegin, dateEnd = history.GetHistoryDate()
	} else {
		dateBegin = date
		dateEnd = time.Date(dateBegin.Year(), dateBegin.Month()+1, 0, 23, 59, 59, 999999999, dateBegin.Location()).UTC()
	}

	dates := make([]time.Time, dateEnd.Day())
	for index := range dates {
		dates[index] = dateBegin.AddDate(0, 0, index)
		date := dates[index].Format("2006-01-02")
		header = append(header, newCell(date + " - Cost").addStyle(textCenter, textBold, backgroundGrey))
		if len(header) > 2 {
			header = append(header, newCell(date + " - Variation").addStyle(textCenter, textBold, backgroundGrey),)
		}
	}

	data = append(data, header)

	for _, aa := range aas {
		report, err := diff.TaskDiffData(ctx, aa, date)
		if err != nil {
			logger.Error("An error occured while generating a cost differentiator report", map[string]interface{}{
				"error": err,
				"account": aa,
			})
			return data, err
		}
		for product, values := range report {
			row := make([]cell, 0)
			row = append(row, newCell(aa.AwsIdentity).addStyle(backgroundLightGrey))
			row = append(row, newCell(product).addStyle(backgroundLightGrey))
			formattedValues, err := formatCostDiff(values)
			if err != nil {
				logger.Error("An error occured while parsing timestamp", map[string]interface{}{
					"error": err,
					"account": aa,
					"values": values,
				})
				return data, err
			}
			for _, date := range dates {
				value := getValueForDate(formattedValues, date)
				if value != nil {
					row = append(row, newCell(value.Cost))
				} else {
					row = append(row, newCell("N/A"))
				}
				if len(row) > 4 {
					if value != nil {
						cell := newCell(value.PercentVariation / 100)
						if value.PercentVariation < 0 {
							cell.addStyle(backgroundGreen)
						} else if value.PercentVariation > 0 {
							cell.addStyle(backgroundRed)
						}
						row = append(row, cell)
					} else {
						row = append(row, newCell("N/A"))
					}
				}
			}
			data = append(data, row)
		}
	}
	return
}
