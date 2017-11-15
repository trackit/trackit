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

package s3

import (
	"database/sql"

	"github.com/trackit/trackit2/aws"
	"github.com/trackit/trackit2/models"
)

// BillRepository is a location where the server may look for bill objects.
type BillRepository struct {
	Bucket string
	Prefix string
}

func CreateBillRepository(aa aws.AwsAccount, br BillRepository, tx *sql.Tx) (BillRepository, error) {
	dbbr := models.AwsBillRepository{
		Prefix:       br.Prefix,
		Bucket:       br.Bucket,
		AwsAccountID: aa.Id,
	}
	var out BillRepository
	err := dbbr.Insert(tx)
	if err == nil {
		out = billRepoFromDbBillRepo(dbbr)
	}
	return out, err
}

func GetBillRepositoriesForAwsAccount(aa aws.AwsAccount, tx *sql.Tx) ([]BillRepository, error) {
	dbAwsBillRepositories, err := models.AwsBillRepositoriesByAwsAccountID(tx, aa.Id)
	if err == nil {
		out := make([]BillRepository, len(dbAwsBillRepositories))
		for i := range out {
			out[i] = billRepoFromDbBillRepo(*dbAwsBillRepositories[i])
		}
		return out, nil
	} else {
		return nil, err
	}
}

func billRepoFromDbBillRepo(dbBillRepo models.AwsBillRepository) BillRepository {
	return BillRepository{
		Bucket: dbBillRepo.Bucket,
		Prefix: dbBillRepo.Prefix,
	}
}
