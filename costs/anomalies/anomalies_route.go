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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/anomaliesDetection"
	"github.com/trackit/trackit-server/cache"
	"github.com/trackit/trackit-server/config"
	"github.com/trackit/trackit-server/costs/anomalies/anomalyFilters"
	"github.com/trackit/trackit-server/costs/anomalies/anomalyType"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/errors"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/models"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

type (
	// esProductAnomalyTypedResult is used to store the raw ElasticSearch response.
	esProductAnomalyTypedResult struct {
		Id        string `json:"-"`
		Account   string `json:"account"`
		Date      string `json:"date"`
		Product   string `json:"product"`
		Abnormal  bool   `json:"abnormal"`
		Recurrent bool   `json:"recurrent"`
		Cost      struct {
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
			cache.UsersCache{},
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
func makeElasticSearchRequest(ctx context.Context, parsedParams anomalyType.AnomalyEsQueryParams) (*elastic.SearchResult, int, error) {
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
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details != nil && cast.Details.Type == "search_phase_execution_exception" {
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

// getSnoozedAnomalies get all snoozed anomalies in database.
func getSnoozedAnomalies(userId int, tx *sql.Tx) (map[string]bool, error) {
	if snoozedAnomalies, err := models.AnomalySnoozingsByUserID(tx, userId); err != nil {
		return nil, err
	} else {
		res := make(map[string]bool)
		for _, snoozedAnomaly := range snoozedAnomalies {
			res[snoozedAnomaly.AnomalyID] = true
		}
		return res, nil
	}
}

// getAnomalyLevel get anomaly level depending on their cost.
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

func formatAnomaliesData(raw *elastic.SearchResult, snoozedAnomalies map[string]bool, ctx context.Context) (anomalyType.AnomaliesDetectionResponse, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	res := make(anomalyType.AnomaliesDetectionResponse)
	for i := range raw.Hits.Hits {
		var typedDocument esProductAnomalyTypedResult
		typedDocument.Id = raw.Hits.Hits[i].Id
		if err := json.Unmarshal(*raw.Hits.Hits[i].Source, &typedDocument); err != nil {
			logger.Error("Failed to parse elasticsearch document.", err.Error())
			return nil, errors.GetErrorMessage(ctx, err)
		}
		if _, ok := res[typedDocument.Account]; !ok {
			res[typedDocument.Account] = make(anomalyType.ProductAnomalies)
		}
		if _, ok := res[typedDocument.Account][typedDocument.Product]; !ok {
			res[typedDocument.Account][typedDocument.Product] = make([]anomalyType.ProductAnomaly, 0)
		}
		level, prettyLevel := getAnomalyLevel(typedDocument)
		if date, err := time.Parse("2006-01-02T15:04:05.000Z", typedDocument.Date); err == nil {
			res[typedDocument.Account][typedDocument.Product] = append(res[typedDocument.Account][typedDocument.Product], anomalyType.ProductAnomaly{
				Id:          typedDocument.Id,
				Date:        date,
				Cost:        typedDocument.Cost.Value,
				UpperBand:   typedDocument.Cost.MaxExpected,
				Abnormal:    typedDocument.Abnormal,
				Recurrent:   typedDocument.Recurrent,
				Filtered:    false,
				Snoozed:     snoozedAnomalies[typedDocument.Id],
				Level:       level,
				PrettyLevel: prettyLevel,
			})
		}
	}
	return res, nil
}

// removeNormalProduct removes product without any anomalies.
func removeNormalProduct(res anomalyType.AnomaliesDetectionResponse) anomalyType.AnomaliesDetectionResponse {
	for account := range res {
		keyToDelete := make([]string, 0)
	pLoop:
		for product := range res[account] {
			for _, an := range res[account][product] {
				if an.Abnormal {
					continue pLoop
				}
			}
			keyToDelete = append(keyToDelete, product)
		}
		for _, key := range keyToDelete {
			delete(res[account], key)
		}
	}
	return res
}

// applyFilters will apply all filters to the response.
func applyFilters(res anomalyType.AnomaliesDetectionResponse, user users.User, ctx context.Context, tx *sql.Tx) anomalyType.AnomaliesDetectionResponse {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	if dbUser, err := models.UserByID(tx, user.Id); err != nil {
		l.Error("Failed to get user with id", map[string]interface{}{
			"userId": user.Id,
			"error":  err.Error(),
		})
	} else {
		filters := FiltersBody{anomalyType.Filters{}}
		if dbUser.AnomaliesFilters != nil {
			if err := json.Unmarshal(dbUser.AnomaliesFilters, &filters.Filters); err != nil {
				l.Error("Failed to unmarshal anomalies filters", map[string]interface{}{
					"userId": user.Id,
					"error":  err.Error(),
				})
			} else {
				res = anomalyFilters.Apply(filters.Filters, res)
			}
		}
	}
	return res
}

// getAnomaliesData checks the request and returns AnomaliesData.
func getAnomaliesData(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := anomalyType.AnomalyEsQueryParams{
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
	snoozedAnomalies, err := getSnoozedAnomalies(user.Id, tx)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	res, err := formatAnomaliesData(raw, snoozedAnomalies, request.Context())
	if err != nil {
		return http.StatusInternalServerError, err
	}
	res = removeNormalProduct(res)
	return http.StatusOK, applyFilters(res, user, request.Context(), tx)
}
