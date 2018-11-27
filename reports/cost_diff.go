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
	"sort"
	"strings"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/costs/diff"
)

type costDiffValue struct {
	Cost      float64
	Variation float64
}

type costDiffProduct map[string]costDiffValue

func isInList(dateList []string, date string) bool {
	for _, item := range dateList {
		if item == date {
			return true
		}
	}
	return false
}

func getDates(rawData map[string]costDiffProduct) (dateList []string) {
	dateList = make([]string, 0)
	for _, product := range rawData {
		for date := range product {
			if !isInList(dateList, date) {
				dateList = append(dateList, date)
			}
		}
	}
	sort.Strings(dateList)
	return
}

func formatCostDiff(data []diff.PricePoint) (values costDiffProduct) {
	values = make(costDiffProduct)
	for _, key := range data {
		var date string
		if splittedDate := strings.Split(key.Date, "T"); len(splittedDate) > 0 {
			date = splittedDate[0]
		} else {
			date = key.Date
		}
		values[date] = costDiffValue{
			Cost: key.Cost,
			Variation: key.PercentVariation / 100,
		}
	}
	return
}

func getCostDiff(ctx context.Context, aa aws.AwsAccount, tx *sql.Tx) (data [][]cell, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Getting Cost Differentiator Report for account", map[string]interface{}{
		"account": aa,
	})

	data = make([][]cell, 0)
	header := make([]cell, 0)
	header = append(header, newCell("Product").addStyle(textCenter, textBold, backgroundGrey))

	rawData := make(map[string]costDiffProduct)

	report, err := diff.TaskDiffData(ctx, aa)
	if err != nil {
		logger.Error("An error occured while generating a cost differentiator report", err)
		return
	}
	for product, data := range report {
		rawData[product] = formatCostDiff(data)
	}

	dates := getDates(rawData)
	for _, date := range dates {
		header = append(header, newCell(date + " - Cost").addStyle(textCenter, textBold, backgroundGrey))
		if len(header) > 2 {
			header = append(header, newCell(date + " - Variation").addStyle(textCenter, textBold, backgroundGrey),)
		}
	}
	data = append(data, header)
	for product, value := range rawData {
		row := make([]cell, 0)
		row = append(row, newCell(product).addStyle(backgroundLightGrey))
		for _, date := range dates {
			if cost, ok := value[date] ; ok {
				row = append(row, newCell(cost.Cost))
				if len(row) > 2 {
					cell := newCell(cost.Variation)
					if cost.Variation < 0 {
						cell.addStyle(backgroundGreen)
					} else if cost.Variation > 0 {
						cell.addStyle(backgroundRed)
					}
					row = append(row, cell)
				}
			} else {
				row = append(row, newCell("N/A"))
				row = append(row, newCell("N/A"))
			}
		}
		data = append(data, row)
	}
	return
}
