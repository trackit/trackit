//   Copyright 2017-2018 MSolution.IO
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

package tags

import (
	"context"
	"errors"
	"net/http"
	"encoding/json"

	"gopkg.in/olivere/elastic.v5"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/users"
	"github.com/trackit/trackit-server/es"
)

type (
	esTagsKeysResult struct {
		Tags struct {
			Buckets []struct {
				Key string `json:"key"`
			} `json:"buckets"`
		}
	}

	TagsKeys []string
)

func getTagsKeysWithParsedParams(ctx context.Context, params tagsKeysQueryParams, user users.User) (int, interface{}){
	var typedDocument esTagsKeysResult
	var response = TagsKeys{}
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	res, returnCode, err := makeElasticSearchRequestForTagsKeys(ctx, params, user, es.Client)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, response
		}
		return returnCode, errors.New("Internal server error")
	}
	err = json.Unmarshal(*res.Aggregations["tags"], &typedDocument.Tags)
	if err != nil {
		l.Error("Error while unmarshaling", err)
		return http.StatusInternalServerError, errors.New("Internal server error")
	}
	for _, tag := range typedDocument.Tags.Buckets {
		response = append(response, tag.Key)
	}
	return http.StatusOK, response
}

func makeElasticSearchRequestForTagsKeys(ctx context.Context, params tagsKeysQueryParams, user users.User, client *elastic.Client) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := getTagsKeysQuery(params)
	index := es.IndexNameForUser(user, "lineitems")
	script := elastic.NewScript("if (params._source.containsKey(\"tags\")) {  params._source[\"tags\"].keySet(); }")
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("tags", elastic.NewTermsAggregation().Script(script))
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists : " + index, err)
			return nil, http.StatusOK, err
		}
		l.Error("Query execution failed : " + err.Error(), nil)
		return nil, http.StatusInternalServerError, err
	}
	return res, http.StatusOK, nil
}

func getTagsKeysQuery(params tagsKeysQueryParams) (*elastic.BoolQuery) {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilter(params.AccountList))
	}
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(params.DateBegin).To(params.DateEnd))
	return query
}
