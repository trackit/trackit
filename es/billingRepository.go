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

package es

import (
	"context"

	"github.com/olivere/elastic"

	"github.com/trackit/trackit/es/indexes/lineItems"
)

// CleanByBillRepositoryId removes every bills information of a specific bill repository
func CleanByBillRepositoryId(ctx context.Context, aaUId, brId int) error {
	index := IndexNameForUserId(aaUId, lineItems.Model.IndexSuffix)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("billRepositoryId", brId))
	_, err := elastic.NewDeleteByQueryService(Client).WaitForCompletion(false).Index(index).Query(query).Do(ctx)
	return err
}

// CleanCurrentMonthBillByBillRepositoryId removes incomplete bills of a specific bill repository (invoiceId == "" when incomplete)
func CleanCurrentMonthBillByBillRepositoryId(ctx context.Context, aaUId, brId int) error {
	index := IndexNameForUserId(aaUId, lineItems.Model.IndexSuffix)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("billRepositoryId", brId), elastic.NewTermQuery("invoiceId", ""))
	_, err := elastic.NewDeleteByQueryService(Client).WaitForCompletion(false).Index(index).Query(query).Do(ctx)
	return err
}
