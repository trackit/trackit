//   Copyright 2017 MSolution.IO
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

package rds

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/rds"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getRDSReport).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.Documentation{
				Summary:     "get the latest RDS report",
				Description: "Responds with the latest RDS report for the account specified in the request",
			},
			routes.QueryArgs{routes.AwsAccountQueryArg},
		),
	}.H().Register("/rds")
}

func makeElasticSearchRequest(ctx context.Context, account string, user users.User) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUser(user, rds.IndexPrefixRDSReport)
	accountFormatted := make([]interface{}, 1)
	accountFormatted[0] = account
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermsQuery("account", accountFormatted...))
	res, err := es.Client.Search().Index(index).Query(query).Sort("reportDate", false).Size(1).Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists : "+index, err)
			return nil, http.StatusOK, err
		}
		l.Error("Query execution failed : "+err.Error(), nil)
		return nil, http.StatusInternalServerError, fmt.Errorf("could not execute the ElasticSearch query")
	}
	if len(res.Hits.Hits) == 0 {
		return nil, http.StatusNotFound, fmt.Errorf("no report found")
	}
	return res, http.StatusOK, nil
}

func getRDSReport(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	account := ""
	if a[routes.AwsAccountQueryArg] != nil {
		account = a[routes.AwsAccountQueryArg].(string)
	}
	if err := aws.ValidateAwsAccounts([]string{account}); err != nil {
		return http.StatusBadRequest, err
	}
	searchResult, returnCode, err := makeElasticSearchRequest(request.Context(), account, user)
	if err != nil {
		return returnCode, err
	}
	return http.StatusOK, (*searchResult.Hits.Hits[0]).Source
}
