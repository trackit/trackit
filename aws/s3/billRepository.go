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
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/aws"
	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/models"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getBillRepository).With(
			routes.Documentation{
				Summary:     "get aws account's bill repositories",
				Description: "Gets the list of bill repositories for an AWS account.",
			},
		),
		http.MethodPost: routes.H(postBillRepository).With(
			routes.RequestContentType{"application/json"},
			routes.Documentation{
				Summary:     "add a new bill repository to an aws account",
				Description: "Adds a bill repository to an AWS account.",
			},
		),
	}.H().With(
		db.RequestTransaction{db.Db},
		users.RequireAuthenticatedUser{},
		routes.RequiredQueryArgs{aws.AwsAccountQueryArg},
		aws.RequireAwsAccount{},
		routes.Documentation{
			Summary:     "interact with aws account's bill repositories",
			Description: "A bill repository is an S3 location (bucket+prefix) where Cost And Usage Reports can be found.",
		},
	).Register("/aws/billrepository")
}

const (
	reportUpdateInterval        = 12 * time.Hour
	reportUpdateVariationAfter  = 6 * time.Hour
	reportUpdateVariationBefore = 2 * time.Hour
)

// BillRepository is a location where the server may look for bill objects.
type BillRepository struct {
	Id                 int       `json:"id"`
	AwsAccountId       int       `json:"awsAccountId"`
	Bucket             string    `json:"bucket"`
	Prefix             string    `json:"prefix"`
	LastImportedPeriod time.Time `json:"lastImportedPeriod"`
	NextUpdate         time.Time `json:"nextUpdate"`
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
		Id:     dbBillRepo.ID,
		Bucket: dbBillRepo.Bucket,
		Prefix: dbBillRepo.Prefix,
	}
}

func postBillRepository(r *http.Request, a routes.Arguments) (int, interface{}) {
	var body BillRepository
	err := decodeRequestBody(r, &body)
	if err == nil {
		err = isBillRepositoryValid(body)
	}
	if err == nil {
		tx := a[db.Transaction].(*sql.Tx)
		aa := a[aws.AwsAccountSelection].(aws.AwsAccount)
		return postBillRepositoryWithValidBody(r, tx, aa, body)
	} else {
		return http.StatusBadRequest, errors.New(fmt.Sprintf("Body is invalid (%s).", err.Error()))
	}
}

// decodeRequestBody decodes a JSON request body and returns nil in case it
// could do so.
func decodeRequestBody(request *http.Request, structuredBody interface{}) error {
	return json.NewDecoder(request.Body).Decode(structuredBody)
}

func postBillRepositoryWithValidBody(r *http.Request, tx *sql.Tx, aa aws.AwsAccount, body BillRepository) (int, interface{}) {
	br, err := CreateBillRepository(aa, body, tx)
	if err == nil {
		go UpdateReport(context.Background(), aa, br)
		return http.StatusOK, br
	} else {
		l := jsonlog.LoggerFromContextOrDefault(r.Context())
		l.Error("Failed to create bill repository.", map[string]interface{}{
			"billRepository": br,
			"error":          err.Error(),
		})
		return http.StatusInternalServerError, errors.New("failed to create bill repository")
	}

}

const noTwoDotsInBucketNameRegex = `^[a-z-](?:[a-z0-9.-]?[a-z0-9-])+$`

var noTwoDotsInBucketName = regexp.MustCompile(noTwoDotsInBucketNameRegex)

func isBillRepositoryValid(br BillRepository) error {
	if err := isBucketNameValid(br.Bucket); err != nil {
		return err
	} else if err := isPrefixValid(br.Prefix); err != nil {
		return err
	} else {
		return nil
	}
}

func isBucketNameValid(bn string) error {
	if l := len(bn); l < 3 {
		return errors.New("bucket name shall be no shorter than 3 chars")
	} else if l > 63 {
		return errors.New("bucket name shall be no longer than 63 chars")
	} else if !noTwoDotsInBucketName.MatchString(bn) {
		return errors.New(fmt.Sprintf("bucket name shall satisfy the regexp /%s/", noTwoDotsInBucketNameRegex))
	} else {
		return nil
	}
}

func isPrefixValid(p string) error {
	l := len([]byte(p))
	if l > 1024 {
		return errors.New("key prefix shall be no longer than 1024 chars")
	} else {
		return nil
	}
}

func getBillRepository(r *http.Request, a routes.Arguments) (int, interface{}) {
	aa := a[aws.AwsAccountSelection].(aws.AwsAccount)
	tx := a[db.Transaction].(*sql.Tx)
	if brs, err := GetBillRepositoriesForAwsAccount(aa, tx); err == nil {
		return http.StatusOK, brs
	} else {
		l := jsonlog.LoggerFromContextOrDefault(r.Context())
		l.Error("Failed to get aws account's bill repositories.", map[string]interface{}{
			"user":       a[users.AuthenticatedUser].(users.User),
			"awsAccount": aa,
			"error":      err.Error(),
		})
		return http.StatusInternalServerError, errors.New("Failed to retrieve bill repositories.")
	}

}

/*
func AwsBillRepositoriesWithDueUpdate(ctx context.Context, db models.XODB) ([]BillRepository, error) {
	dbRepos, err := models.AwsBillRepositoriesWithDueUpdate(db)
	if err != nil {

	}
}
*/
