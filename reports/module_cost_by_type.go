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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/olivere/elastic"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/users"
)

const costByTypeSheetName = "Cost by type"
const maxAggregationSize = 0x7FFFFFFF

var costByTypeModule = module{
	Name:          "Cost by type",
	SheetName:     costByTypeSheetName,
	ErrorName:     "costByTypeError",
	GenerateSheet: generateCostByTypeSheet,
}

type QueryParams struct {
	AccountList []string  `json:"awsAccounts"`
	IndexList   []string  `json:"indexes"`
	DateBegin   time.Time `json:"begin"`
	DateEnd     time.Time `json:"end"`
}

type esCostByTypesResult struct {
	Buckets []struct {
		Key      string `json:"key"`
		Services struct {
			Buckets []struct {
				Key       string `json:"key"`
				Resources struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
						Cost     struct {
							Value float32 `json:"value"`
						} `json:"cost"`
					} `json:"buckets"`
				} `json:"resources"`
			} `json:"buckets"`
		} `json:"services"`
	} `json:"buckets"`
}

// generateCostByTypeSheet will generate a sheet with Cost by type
// It will get data for given AWS account and for a given date
func generateCostByTypeSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx, file *excelize.File) (err error) {
	if date.IsZero() {
		date, _ = history.GetHistoryDate()
	}
	return costByTypeGenerateSheet(ctx, aas, date, file, tx)
}

func costByTypeGenerateSheet(ctx context.Context, aas []aws.AwsAccount, date time.Time, file *excelize.File, tx *sql.Tx) (err error) {
	data, err := costByTypeGetData(ctx, aas, date, tx)

	if err == nil {
		costByTypeInsertDataInSheet(file, costByTypeSheetName, data)
	}
	return err
}

func createQueryAccountFilter(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("usageAccountId", accountListFormatted...)
}

func getQueryWithParams(params QueryParams) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilter(params.AccountList))
	}
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(params.DateBegin).To(params.DateEnd))
	return query
}

func makeElasticSearchRequestForCostType(ctx context.Context, params QueryParams,
	client *elastic.Client) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := getQueryWithParams(params)
	index := strings.Join(params.IndexList, ",")
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("costGroups", elastic.NewTermsAggregation().Field("costGroup").Size(maxAggregationSize).
		SubAggregation("services", elastic.NewTermsAggregation().Field("serviceCode").Size(maxAggregationSize).
			SubAggregation("resources", elastic.NewTermsAggregation().Field("usageType").Size(maxAggregationSize).
				SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost")).Size(maxAggregationSize))))
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, http.StatusOK, err
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details != nil && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, err
	}
	return res, http.StatusOK, nil
}

func costByTypeGetData(ctx context.Context, aas []aws.AwsAccount, date time.Time, tx *sql.Tx) (data esCostByTypesResult, err error) {
	var doc esCostByTypesResult

	client := es.Client
	identities := getAwsIdentities(aas)

	user, err := users.GetUserWithId(tx, aas[0].UserId)
	if err != nil {
		return doc, err
	}

	queryParams := QueryParams{
		AccountList: []string{},
		IndexList:   []string{},
		DateBegin:   date,
		DateEnd:     time.Date(date.Year(), date.Month()+1, 0, 23, 59, 59, 999999999, date.Location()).UTC(),
	}

	accountsAndIndexes, _, err := es.GetAccountsAndIndexes(identities, user, tx, s3.IndexPrefixLineItem)
	if err != nil {
		return doc, err
	}
	queryParams.AccountList = accountsAndIndexes.Accounts
	queryParams.IndexList = accountsAndIndexes.Indexes

	res, _, err := makeElasticSearchRequestForCostType(ctx, queryParams, client)
	if err != nil {
		return doc, err
	}
	err = json.Unmarshal(*res.Aggregations["costGroups"], &doc)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func costByTypeGenerateSheetHeader(file *excelize.File, name string) {
	header := cells{
		newCell("Network", "A1").mergeTo("D1"),
		newCell("Storage", "F1").mergeTo("I1"),
		newCell("Compute", "K1").mergeTo("N1"),

		newCell("", "A2").mergeTo("D2"),
		newCell("", "F2").mergeTo("I2"),
		newCell("", "K2").mergeTo("N2"),

		newCell("Service", "A3"),
		newCell("Type", "B3"),
		newCell("Amount", "C3"),
		newCell("Cost", "D3"),

		newCell("Service", "F3"),
		newCell("Type", "G3"),
		newCell("Amount", "H3"),
		newCell("Cost", "I3"),

		newCell("Service", "K3"),
		newCell("Type", "L3"),
		newCell("Amount", "M3"),
		newCell("Cost", "N3"),
	}
	header.addStyles("borders", "bold", "centerText").setValues(file, name)

	rows := rowsHeight{
		newRowHeight(2, 40),
	}
	rows.setValues(file, name)

	columns := columnsWidth{
		newColumnWidth("A", 30),
		newColumnWidth("F", 30),
		newColumnWidth("K", 30),

		newColumnWidth("B", 45),
		newColumnWidth("G", 45),
		newColumnWidth("L", 45),

		newColumnWidth("C", 15),
		newColumnWidth("H", 15),
		newColumnWidth("M", 15),

		newColumnWidth("D", 15),
		newColumnWidth("I", 15),
		newColumnWidth("N", 15),
	}
	columns.setValues(file, name)
}

func costByTypeInsertDataInSheet(file *excelize.File, name string, data esCostByTypesResult) {
	file.NewSheet(name)
	costByTypeGenerateSheetHeader(file, name)

	xOffset := 0
	yOffsetService := 4
	yOffsetResource := 4

loop:
	for _, costType := range data.Buckets {
		switch costType.Key {
		case "network":
			xOffset = 0
			break
		case "storage":
			xOffset = 5
			break
		case "compute":
			xOffset = 10
			break
		default:
			continue loop
		}

		yOffsetService = 4
		yOffsetResource = 4

		serviceCol := string(int('A') + xOffset)
		resourceCol := string(int('A') + xOffset + 1)
		amountCol := string(int('A') + xOffset + 2)
		costCol := string(int('A') + xOffset + 3)

		totalCol := string(int('A') + xOffset)
		endTotalCol := string(int('A') + xOffset + 3)

		for _, service := range costType.Services.Buckets {
			serviceHeight := len(service.Resources.Buckets)
			var cell cell

			if serviceHeight > 1 {
				cell = newCell(service.Key, serviceCol+strconv.Itoa(yOffsetService)).mergeTo(serviceCol + strconv.Itoa(yOffsetService+serviceHeight-1))
			} else {
				cell = newCell(service.Key, serviceCol+strconv.Itoa(yOffsetService))
			}
			cell.addStyles("borders", "centerText").setValue(file, name)
			yOffsetService += serviceHeight

			for _, resource := range service.Resources.Buckets {
				cells := cells{
					newCell(resource.Key, resourceCol+strconv.Itoa(yOffsetResource)),
					newCell(resource.DocCount, amountCol+strconv.Itoa(yOffsetResource)),
					newCell(resource.Cost.Value, costCol+strconv.Itoa(yOffsetResource)).addStyles("price"),
				}
				cells.addStyles("borders", "centerText").setValues(file, name)
				yOffsetResource += 1
			}
		}
		if yOffsetResource != 4 {
			newFormula("SUM("+costCol+"4:"+costCol+strconv.Itoa(yOffsetResource-1)+")", totalCol+"2").mergeTo(endTotalCol+"2").addStyles("borders", "centerText", "bold", "price").setValue(file, name)
		}
	}
}
