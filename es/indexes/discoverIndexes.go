//   Copyright 2020 MSolution.IO
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

package indexes

import (
	"context"
	"strings"

	"github.com/trackit/trackit/es/indexes/common"

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/models"
)

func discoverIndexes(ctx context.Context) error {
	indexNames, err := es.Client.IndexNames()
	if err != nil {
		return err
	}

	for _, indexName := range indexNames {
		indexVersioningData := getIndexVersioningData(indexName)
		if indexVersionned(indexName) || indexVersioningData == nil {
			continue
		}
		dbObject := models.EsVersioning{
			CurrentVersion: indexVersioningData.Version,
			TemplateName:   indexVersioningData.Name,
			IndexName:      indexName,
		}
		err = dbObject.Insert(db.Db)
		if err != nil {
			return err
		}
	}

	return nil
}

func indexVersionned(indexName string) bool {
	res, _ := models.EsVersioningByIndexName(db.Db, indexName)
	return res != nil
}

func getIndexVersioningData(indexName string) *common.VersioningData {
	for _, template := range versioningData {
		if strings.HasSuffix(indexName, template.IndexSuffix) {
			return &template
		}
	}
	return nil
}
