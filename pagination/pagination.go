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

package pagination

import (
	"fmt"
	"github.com/trackit/trackit/routes"
)

type Pagination struct {
	Page          int   `json:"page"`
	Elements      int   `json:"elements"`
	TotalElements int   `json:"total_elements"`
}

const (
	DefaultPageForPagination = 1
	DefaultTotalElements     = 50
	MaxAggregationSize       = 0x7FFFFFFF
)

func (page Pagination) GetFromValue() int {
	if page.Page == 1 {
		return 0
	}
	//es.Client.Scroll("*-lambda-reports").
	return page.Elements * (page.Page - 1)
}

func (page Pagination) setDefaultAggregationSize() Pagination {
	page.Elements = MaxAggregationSize
	page.Page = DefaultPageForPagination
	return page
}

func StoreTotalHits(page *Pagination, hits int64, formatted int) {
	page.TotalElements = int(hits)
	fmt.Printf("Hits: %v & %v\n", hits, formatted)
	if formatted < page.TotalElements {
		page.TotalElements = formatted
	}
}

func NewPagination(arguments routes.Arguments) Pagination {
	var page Pagination
	if arguments == nil {
		return page.setDefaultAggregationSize()
	}
	if arguments[routes.PaginationPageQueryArg] != nil {
		value := arguments[routes.PaginationPageQueryArg].(int)
		page.Page = value
	}
	if arguments[routes.PaginationNumberElementsQueryArg] != nil {
		value := arguments[routes.PaginationNumberElementsQueryArg].(int)
		page.Elements = value
	}
	if page.Page <= 0 {
		page.Page = DefaultPageForPagination
	}
	if page.Elements <= 0 {
		page.Elements = DefaultTotalElements
	}
	if page.GetFromValue() > MaxAggregationSize {
		page.Page = DefaultPageForPagination
		page.Elements = DefaultTotalElements
	}
	return page
}
