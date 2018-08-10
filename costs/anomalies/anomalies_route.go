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
	"fmt"
	"net/http"
	"time"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"

	"gopkg.in/olivere/elastic.v5"
)

// esQueryParams will store the parsed query params
type esQueryParams struct {
	dateBegin         time.Time
	dateEnd           time.Time
	accountList       []string
}

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

// makeElasticSearchRequest prepares and run the request to retrieve the cost anomalies.
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequest(ctx context.Context, parsedParams esQueryParams, user users.User) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUser(user, "lineitems")
	searchService := GetElasticSearchParams(
		parsedParams.accountList,
		parsedParams.dateBegin,
		parsedParams.dateEnd,
		"12h",
		es.Client,
		index,
	)
	res, err := searchService.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists : "+index, err)
			return nil, http.StatusOK, err
		}
		l.Error("Query execution failed : "+err.Error(), nil)
		return nil, http.StatusInternalServerError, fmt.Errorf("could not execute the ElasticSearch query")
	}
	return res, http.StatusOK, nil
}

// getAnomaliesData returns the cost anomalies based on the query params, in JSON format.
func getAnomaliesData(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := esQueryParams{
		accountList:       []string{},
		dateBegin:         a[anomalyQueryArgs[1]].(time.Time),
		dateEnd:           a[anomalyQueryArgs[2]].(time.Time).Add(time.Hour*time.Duration(23) + time.Minute*time.Duration(59) + time.Second*time.Duration(59)),
	}
	if a[anomalyQueryArgs[0]] != nil {
		parsedParams.accountList = a[anomalyQueryArgs[0]].([]string)
	}
	if err := aws.ValidateAwsAccounts(parsedParams.accountList); err != nil {
		return http.StatusBadRequest, err
	}
	sr, returnCode, err := makeElasticSearchRequest(request.Context(), parsedParams, user)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, nil
		} else {
			return returnCode, err
		}
	}
	res, err := prepareAnomalyData(request.Context(), sr)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, res
}
