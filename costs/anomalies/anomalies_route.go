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

package anomalies

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"encoding/json"
	"github.com/trackit/trackit-server/anomaliesDetection"
	"github.com/trackit/trackit-server/config"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
	"strconv"
)

type (
	// AnomalyEsQueryParams will store the parsed query params
	AnomalyEsQueryParams struct {
		DateBegin   time.Time
		DateEnd     time.Time
		AccountList []string
		IndexList   []string
		AnomalyType string
	}

	// productAnomaly represents one anomaly returned.
	productAnomaly struct {
		Date        time.Time `json:"date"`
		Cost        float64   `json:"cost"`
		UpperBand   float64   `json:"upper_band"`
		Abnormal    bool      `json:"abnormal"`
		Level       int       `json:"level"`
		PrettyLevel string    `json:"pretty_level"`
	}

	// productAnomalies is used to respond to the request.
	// Key is a product name.
	productAnomalies map[string][]productAnomaly

	// anomaliesDetectionResponse is used to respond to the request.
	// Key is an AWS Account Identity.
	anomaliesDetectionResponse map[string]productAnomalies

	// esProductAnomalyTypedResult is used to store the raw ElasticSearch response.
	esProductAnomalyTypedResult struct {
		Account  string `json:"account"`
		Date     string `json:"date"`
		Product  string `json:"product"`
		Abnormal bool   `json:"abnormal"`
		Cost     struct {
			Value       float64 `json:"value"`
			MaxExpected float64 `json:"maxExpected"`
		} `json:"cost"`
	}
)

// anomalyQueryArgs allows to get required queryArgs params
var anomalyQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
	routes.DateBeginQueryArg,
	routes.DateEndQueryArg,
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getAnomaliesData).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(anomalyQueryArgs),
			routes.Documentation{
				Summary:     "get the cost anomalies",
				Description: "Responds with the cost anomalies based on the query args passed to it",
			},
		),
	}.H().Register("/costs/anomalies")
}

// makeElasticSearchRequest prepares and run the request.
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, parsedParams AnomalyEsQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.IndexList, ",")
	searchService := getElasticSearchParams(
		parsedParams.AccountList,
		parsedParams.DateBegin,
		parsedParams.DateEnd,
		es.Client,
		index,
		parsedParams.AnomalyType,
	)
	res, err := searchService.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, http.StatusOK, errors.GetErrorMessage(ctx, err)
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, http.StatusInternalServerError, errors.GetErrorMessage(ctx, err)
	}
	return res, http.StatusOK, nil
}

func getAnomalyLevel(typedDocument esProductAnomalyTypedResult) (int, string) {
	if !typedDocument.Abnormal {
		return 0, ""
	}
	prettyLevels := strings.Split(config.AnomalyDetectionPrettyLevels, ",")
	percent := (typedDocument.Cost.Value * 100) / typedDocument.Cost.MaxExpected
	levels := strings.Split(config.AnomalyDetectionLevels, ",")
	for i, level := range levels[1:] {
		l, _ := strconv.ParseFloat(level, 64)
		if percent < l {
			return i, prettyLevels[i]
		}
	}
	return len(levels) - 1, prettyLevels[len(levels)-1]
}

func formatAnomaliesData(raw *elastic.SearchResult, ctx context.Context) (anomaliesDetectionResponse, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	res := make(anomaliesDetectionResponse)
	for i := range raw.Hits.Hits {
		var typedDocument esProductAnomalyTypedResult
		if err := json.Unmarshal(*raw.Hits.Hits[i].Source, &typedDocument); err != nil {
			logger.Error("Failed to parse elasticsearch document.", err.Error())
			return nil, errors.GetErrorMessage(ctx, err)
		}
		if _, ok := res[typedDocument.Account]; !ok {
			res[typedDocument.Account] = make(productAnomalies)
		}
		if _, ok := res[typedDocument.Account][typedDocument.Product]; !ok {
			res[typedDocument.Account][typedDocument.Product] = make([]productAnomaly, 0)
		}
		level, prettyLevel := getAnomalyLevel(typedDocument)
		if date, err := time.Parse("2006-01-02T15:04:05.000Z", typedDocument.Date); err == nil {
			res[typedDocument.Account][typedDocument.Product] = append(res[typedDocument.Account][typedDocument.Product], productAnomaly{
				Date:        date,
				Cost:        typedDocument.Cost.Value,
				UpperBand:   typedDocument.Cost.MaxExpected,
				Abnormal:    typedDocument.Abnormal,
				Level:       level,
				PrettyLevel: prettyLevel,
			})
		}
	}
	return res, nil
}

// getAnomaliesData checks the request and returns AnomaliesData.
func getAnomaliesData(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := AnomalyEsQueryParams{
		AccountList: []string{},
		DateBegin:   a[anomalyQueryArgs[1]].(time.Time),
		DateEnd:     a[anomalyQueryArgs[2]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59)),
	}
	if a[anomalyQueryArgs[0]] != nil {
		parsedParams.AccountList = a[anomalyQueryArgs[0]].([]string)
	}
	tx := a[db.Transaction].(*sql.Tx)
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.AccountList, user, tx, anomalies.IndexPrefixAnomaliesDetection)
	if err != nil {
		return returnCode, err
	}
	parsedParams.AccountList = accountsAndIndexes.Accounts
	parsedParams.IndexList = accountsAndIndexes.Indexes
	parsedParams.AnomalyType = anomalies.TypeProductAnomaliesDetection
	raw, returnCode, err := makeElasticSearchRequest(request.Context(), parsedParams)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, err
		} else {
			return http.StatusInternalServerError, err
		}
	}
	res, err := formatAnomaliesData(raw, request.Context())
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, res
}
