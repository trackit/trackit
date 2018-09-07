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

package ec2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/ec2"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

// esQueryParams will store the parsed query params
type esQueryParams struct {
	accountList []string
}

// ec2QueryArgs allows to get required queryArgs params
var ec2QueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
}

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getEc2Instances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(ec2QueryArgs),
			routes.Documentation{
				Summary:     "get the list of EC2 instances",
				Description: "Responds with the list of EC2 instances based on the queryparams passed to it",
			},
		),
	}.H().Register("/ec2")
}

// makeElasticSearchRequests prepares and run the request to retrieve the latest reports
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchEc2Request(ctx context.Context, parsedParams esQueryParams, user users.User) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := es.IndexNameForUser(user, ec2.IndexPrefixEC2Report)
	searchService := GetElasticSearchEc2Params(
		parsedParams.accountList,
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

// getEc2Instances returns the list of EC2 reports based on the query params, in JSON format.
func getEc2Instances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := esQueryParams{
		accountList: []string{},
	}
	if a[ec2QueryArgs[0]] != nil {
		parsedParams.accountList = a[ec2QueryArgs[0]].([]string)
	}
	if err := aws.ValidateAwsAccounts(parsedParams.accountList); err != nil {
		return http.StatusBadRequest, err
	}
	searchResult, returnCode, err := makeElasticSearchEc2Request(request.Context(), parsedParams, user)
	if err != nil {
		return returnCode, err
	}
	res, err := prepareResponse(request.Context(), searchResult, user, parsedParams)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, res
}
