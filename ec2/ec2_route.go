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
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/aws/ec2"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

// ec2QueryParams will store the parsed query params
type ec2QueryParams struct {
	accountList []string
	indexList   []string
}

// ec2HistoryQueryParams will store the parsed query params
type ec2HistoryQueryParams struct {
	accountList []string
	indexList   []string
	date        time.Time
}

// ec2QueryArgs allows to get required queryArgs params
var ec2QueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
}

// ec2QueryArgs allows to get required queryArgs params
var ec2HistoryQueryArgs = []routes.QueryArg{
	routes.AwsAccountsOptionalQueryArg,
	routes.QueryArg{
		Name:        "date",
		Type:        routes.QueryArgDate{},
		Description: "Date with year and month. Format is ISO8601",
		Optional:    false,
	},
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
	routes.MethodMuxer{
		http.MethodGet: routes.H(getEc2HistoryInstances).With(
			db.RequestTransaction{Db: db.Db},
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			routes.QueryArgs(ec2HistoryQueryArgs),
			routes.Documentation{
				Summary:     "get the list of EC2 instances of a month",
				Description: "Responds with the list of EC2 instances of a month based on the queryparams passed to it",
			},
		),
	}.H().Register("/ec2/history")
}

// makeElasticSearchCostRequests prepares and run the request to retrieve the cost per instance
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchCostRequest(ctx context.Context, params ec2QueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(params.indexList, ",")
	searchService := GetElasticSearchCostParams(
		params.accountList,
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

// makeElasticSearchEc2Request prepares and run the request to retrieve the latest reports
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchEc2Request(ctx context.Context, parsedParams ec2QueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
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
	parsedParams := ec2QueryParams{
		accountList: []string{},
	}
	if a[ec2QueryArgs[0]] != nil {
		parsedParams.accountList = a[ec2QueryArgs[0]].([]string)
	}
	tx := a[db.Transaction].(*sql.Tx)
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.accountList, user, tx, ec2.IndexPrefixEC2Report)
	if err != nil {
		return returnCode, err
	}
	parsedParams.accountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	searchResult, returnCode, err := makeElasticSearchEc2Request(request.Context(), parsedParams)
	if err != nil {
		return returnCode, err
	}
	accountsAndIndexes, returnCode, err = es.GetAccountsAndIndexes(parsedParams.accountList, user, tx, es.IndexPrefixLineItems)
	if err != nil {
		return returnCode, err
	}
	parsedParams.accountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	costResult, _, _ := makeElasticSearchCostRequest(request.Context(), parsedParams)
	res, err := prepareResponseEc2(request.Context(), searchResult, costResult)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, res
}

// makeElasticSearchEc2HistoryRequest prepares and run the request to retrieve a month report
// based on the esQueryParams
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchEc2HistoryRequest(ctx context.Context, parsedParams ec2HistoryQueryParams) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	index := strings.Join(parsedParams.indexList, ",")
	searchService := GetElasticSearchEc2HistoryParams(
		parsedParams.accountList,
		parsedParams.date,
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

// getEc2HistoryInstances returns the list of EC2 reports based on the query params, in JSON format.
func getEc2HistoryInstances(request *http.Request, a routes.Arguments) (int, interface{}) {
	user := a[users.AuthenticatedUser].(users.User)
	parsedParams := ec2HistoryQueryParams{
		accountList: []string{},
		date:        a[ec2HistoryQueryArgs[1]].(time.Time),
	}
	if a[ec2QueryArgs[0]] != nil {
		parsedParams.accountList = a[ec2HistoryQueryArgs[0]].([]string)
	}
	tx := a[db.Transaction].(*sql.Tx)
	accountsAndIndexes, returnCode, err := es.GetAccountsAndIndexes(parsedParams.accountList, user, tx, ec2.IndexPrefixEC2Report)
	if err != nil {
		return returnCode, err
	}
	parsedParams.accountList = accountsAndIndexes.Accounts
	parsedParams.indexList = accountsAndIndexes.Indexes
	searchResult, returnCode, err := makeElasticSearchEc2HistoryRequest(request.Context(), parsedParams)
	if err != nil {
		return returnCode, err
	}
	res, err := prepareResponseEc2History(request.Context(), searchResult)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, res
}
