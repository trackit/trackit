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
	esTagsValuesResult struct {
		Tags struct {
			Buckets []struct {
				Key string `json:"key"`
				Cost struct {
					Value float64 `json:"value"`
				} `json:"cost"`
			} `json:"buckets"`
		}
	}

	TagValue struct {
		Tag  string  `json:"tag"`
		Cost float64 `json:"cost"`
	}

	TagsValues struct {
		Key        string     `json:"key"`
		TagsValues []TagValue `json:"values"`
	}

	TagsValuesResponse []TagsValues
)

func getTagsValuesWithParsedParams(ctx context.Context, params tagsValuesQueryParams, user users.User) (int, interface{}){
	var response = TagsValuesResponse{}
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	for i := range params.TagsKey {
		var typedDocument esTagsValuesResult
		res, returnCode, err := makeElasticSearchRequestForTagsValues(ctx, params, user, es.Client, i)
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
		tagsValues := TagsValues{params.TagsKey[i], nil}
		for _, tag := range typedDocument.Tags.Buckets {
			tagsValues.TagsValues = append(tagsValues.TagsValues, TagValue{tag.Key, tag.Cost.Value})
		}
		response = append(response, tagsValues)
	}
	return http.StatusOK, response
}

func makeElasticSearchRequestForTagsValues(ctx context.Context, params tagsValuesQueryParams, user users.User, client *elastic.Client, i int) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := getTagsValuesQuery(params)
	index := es.IndexNameForUser(user, "lineitems")
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("tags", elastic.NewTermsAggregation().Field("tags." + params.TagsKey[i]).
		SubAggregation("cost", elastic.NewSumAggregation().Field("unblendedCost")))
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

func getTagsValuesQuery(params tagsValuesQueryParams) (*elastic.BoolQuery) {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilter(params.AccountList))
	}
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(params.DateBegin).To(params.DateEnd))
	return query
}

func createQueryAccountFilter(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("usageAccountId", accountListFormatted...)
}
