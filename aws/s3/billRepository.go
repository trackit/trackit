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
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/es"
	"github.com/trackit/trackit-server/models"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

func init() {
	routes.MethodMuxer{
		http.MethodGet: routes.H(getBillRepository).With(
			users.RequireAuthenticatedUser{users.ViewerAsParent},
			aws.RequireAwsAccountId{},
			routes.Documentation{
				Summary:     "get aws account's bill repositories",
				Description: "Gets the list of bill repositories for an AWS account.",
			},
		),
		http.MethodPost: routes.H(postBillRepository).With(
			users.RequireAuthenticatedUser{users.ViewerCannot},
			aws.RequireAwsAccountId{},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{postBillRepositoryBody{
				Bucket: "my-bucket",
				Prefix: "bills/",
			}},
			routes.Documentation{
				Summary:     "add a new bill repository to an aws account",
				Description: "Adds a bill repository to an AWS account.",
			},
		),
		http.MethodPatch: routes.H(patchBillRepository).With(
			users.RequireAuthenticatedUser{users.ViewerCannot},
			aws.RequireAwsAccountId{},
			routes.RequestContentType{"application/json"},
			routes.QueryArgs{routes.BillPositoryQueryArg},
			routes.RequestBody{postBillRepositoryBody{
				Bucket: "my-bucket",
				Prefix: "bills/",
			}},
			routes.Documentation{
				Summary:     "add a new bill repository to an aws account",
				Description: "Adds a bill repository to an AWS account.",
			},
		),
		http.MethodDelete: routes.H(deleteBillRepository).With(
			users.RequireAuthenticatedUser{users.ViewerCannot},
			aws.RequireAwsAccountId{},
			routes.RequestContentType{"application/json"},
			routes.QueryArgs{routes.BillPositoryQueryArg},
			routes.Documentation{
				Summary:     "delete a bill repository from an aws account",
				Description: "delete a bill repository from an AWS account.",
			},
		),
	}.H().With(
		db.RequestTransaction{db.Db},
		routes.QueryArgs{routes.AwsAccountIdQueryArg},
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
	Id                   int       `json:"id"`
	AwsAccountId         int       `json:"awsAccountId"`
	Bucket               string    `json:"bucket"`
	Prefix               string    `json:"prefix"`
	Error                string    `json:"error"`
	LastImportedManifest time.Time `json:"lastImportedManifest"`
	NextUpdate           time.Time `json:"nextUpdate"`
}

// CreateBillRepository creates a BillRepository for an AwsAccount. It does
// not perform checks on the repository.
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

// UpdateBillRepository updates a BillRepository in the database
func UpdateBillRepositorySafe(dbBr *models.AwsBillRepository, br BillRepository, tx *sql.Tx) (BillRepository, error) {
	dbBr.Prefix = br.Prefix
	dbBr.Bucket = br.Bucket
	dbBr.AwsAccountID = br.AwsAccountId
	dbBr.NextUpdate = br.NextUpdate
	dbBr.LastImportedManifest = br.LastImportedManifest
	dbBr.Error = br.Error
	var out BillRepository
	err := dbBr.Update(tx)
	if err == nil {
		out = billRepoFromDbBillRepo(*dbBr)
	}
	return out, err
}

// UpdateBillRepository updates a BillRepository in the database. No checks are
// performed.
func UpdateBillRepository(br BillRepository, tx *sql.Tx) error {
	dbAwsBillRepository := dbBillRepoFromBillRepo(br)
	return dbAwsBillRepository.UpdateUnsafe(tx)
}

// UpdateBillRepositoryWithoutContext updates a BillRepository in the database.
// No checks are performed.
func UpdateBillRepositoryWithoutContext(br BillRepository, db models.XODB) error {
	dbAwsBillRepository := dbBillRepoFromBillRepo(br)
	return dbAwsBillRepository.UpdateUnsafe(db)
}

// GetBillRepositoriesForAwsAccount retrieves from the database all the
// BillRepositories for an AwsAccount.
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

// GetBillRepositoryForAwsAccountById gets a BillRepository by its ID, ensuring
// it belongs to the provided AwsAccount.
func GetBillRepositoryForAwsAccountById(aa aws.AwsAccount, brId int, tx *sql.Tx) (BillRepository, error) {
	var out BillRepository
	var err error
	dbAwsBillRepository, err := models.AwsBillRepositoryByID(tx, brId)
	if err == nil {
		out = billRepoFromDbBillRepo(*dbAwsBillRepository)
		if out.AwsAccountId != aa.Id {
			err = errors.New("bill repository does not belong to aws account")
		}
	}
	return out, err
}

// GetAwsBillRepositoriesWithDueUpdate gets all bill repositories in need of an
// update.
func GetAwsBillRepositoriesWithDueUpdate(tx *sql.Tx) ([]BillRepository, error) {
	dbbrs, err := models.AwsBillRepositoriesWithDueUpdate(tx)
	if err != nil {
		return nil, err
	}
	brs := make([]BillRepository, len(dbbrs))
	for i := range dbbrs {
		brs[i] = billRepoFromDbBillRepo(*dbbrs[i])
	}
	return brs, nil
}

func billRepoFromDbBillRepo(dbBillRepo models.AwsBillRepository) BillRepository {
	return BillRepository{
		Id:                   dbBillRepo.ID,
		Bucket:               dbBillRepo.Bucket,
		Prefix:               dbBillRepo.Prefix,
		Error:                dbBillRepo.Error,
		AwsAccountId:         dbBillRepo.AwsAccountID,
		LastImportedManifest: dbBillRepo.LastImportedManifest,
		NextUpdate:           dbBillRepo.NextUpdate,
	}
}

func dbBillRepoFromBillRepo(br BillRepository) models.AwsBillRepository {
	return models.AwsBillRepository{
		ID:                   br.Id,
		Bucket:               br.Bucket,
		Prefix:               br.Prefix,
		Error:                br.Error,
		AwsAccountID:         br.AwsAccountId,
		LastImportedManifest: br.LastImportedManifest,
		NextUpdate:           br.NextUpdate,
	}
}

type postBillRepositoryBody struct {
	Prefix string `json:"prefix" req:""`
	Bucket string `json:"bucket" req:"nonzero"`
}

func postBillRepository(r *http.Request, a routes.Arguments) (int, interface{}) {
	var body postBillRepositoryBody
	routes.MustRequestBody(a, &body)
	if err := isBillRepositoryValid(body); err != nil {
		return http.StatusBadRequest, errors.New(fmt.Sprintf("Body is invalid (%s).", err.Error()))
	}
	aa := a[aws.AwsAccountSelection].(aws.AwsAccount)
	if err := isBillRepositoryAccessible(r.Context(), aa, body); err != nil {
		return http.StatusBadRequest, err
	}
	tx := a[db.Transaction].(*sql.Tx)
	return postBillRepositoryWithValidBody(r, tx, aa, body)
}

func postBillRepositoryWithValidBody(
	r *http.Request,
	tx *sql.Tx,
	aa aws.AwsAccount,
	body postBillRepositoryBody,
) (int, interface{}) {
	br, err := CreateBillRepository(aa, BillRepository{Bucket: body.Bucket, Prefix: body.Prefix}, tx)
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

func patchBillRepository(r *http.Request, a routes.Arguments) (int, interface{}) {
	var body postBillRepositoryBody
	routes.MustRequestBody(a, &body)
	if err := isBillRepositoryValid(body); err != nil {
		return http.StatusBadRequest, errors.New(fmt.Sprintf("Body is invalid (%s).", err.Error()))
	}
	aa := a[aws.AwsAccountSelection].(aws.AwsAccount)
	if err := isBillRepositoryAccessible(r.Context(), aa, body); err != nil {
		return http.StatusBadRequest, err
	}
	tx := a[db.Transaction].(*sql.Tx)
	brId := a[routes.BillPositoryQueryArg].(int)
	return patchBillRepositoryWithValidBody(r, tx, aa, brId, body)
}

func patchBillRepositoryWithValidBody(
	r *http.Request,
	tx *sql.Tx,
	aa aws.AwsAccount,
	brId int,
	body postBillRepositoryBody,
) (int, interface{}) {
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	dbBillingRepo, err := models.AwsBillRepositoryByID(tx, brId)
	if err != nil {
		l.Error("Failed to find bill repository to update.", map[string]interface{}{
			"error": err.Error(),
		})
		return http.StatusNotFound, errors.New("failed to find bill repository to update")
	}
	br, err := UpdateBillRepositorySafe(dbBillingRepo, BillRepository{Id: brId, AwsAccountId: aa.Id, Bucket: body.Bucket, Prefix: body.Prefix}, tx)
	if err == nil {
		go func() {
			err = es.CleanByBillRepositoryId(context.Background(), aa.UserId, br.Id)
			if err != nil {
				l.Error("Failed to clean ES data for bill repository", map[string]interface{}{
					"billRepository": br,
					"error":          err.Error(),
				})
			}
			UpdateReport(context.Background(), aa, br)
		}()
		return http.StatusOK, br
	} else {
		l.Error("Failed to update bill repository.", map[string]interface{}{
			"billRepository": br,
			"error":          err.Error(),
		})
		return http.StatusInternalServerError, errors.New("failed to update bill repository")
	}
}

const noTwoDotsInBucketNameRegex = `^[a-z-](?:[a-z0-9.-]?[a-z0-9-])+$`

var noTwoDotsInBucketName = regexp.MustCompile(noTwoDotsInBucketNameRegex)

func isBillRepositoryValid(br postBillRepositoryBody) error {
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

func isBillRepositoryAccessible(ctx context.Context, aa aws.AwsAccount, body postBillRepositoryBody) error {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	_, _, err := getServiceForRepository(ctx, aa, BillRepository{Bucket: body.Bucket, Prefix: body.Prefix})
	if err != nil {
		l.Warning("Trying to add a bad bill location.", err.Error())
		return errors.New("Couldn't access to this bill location.")
	}
	return nil
}

func DeleteBillRepositoryById(brId int, tx *sql.Tx) error {
	dbBr, err := models.AwsBillRepositoryByID(tx, brId)
	if err == nil {
		return dbBr.Delete(tx)
	} else if err.Error() == "sql: no rows in result set" {
		return errors.New("Failed to retrieve bill repository.")
	}
	return err
}

func deleteBillRepository(r *http.Request, a routes.Arguments) (int, interface{}) {
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	aa := a[aws.AwsAccountSelection].(aws.AwsAccount)
	brId := a[routes.BillPositoryQueryArg].(int)
	tx := a[db.Transaction].(*sql.Tx)
	err := DeleteBillRepositoryById(brId, tx)
	if err == nil {
		go func() {
			err = es.CleanByBillRepositoryId(context.Background(), aa.UserId, brId)
			if err != nil {
				l.Error("Failed to clean ES data for bill repository", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
		return http.StatusOK, nil
	} else if err.Error() == "Failed to retrieve bill repository." {
		l.Error("Failed to delete billing repository.", err.Error())
		return http.StatusNotFound, errors.New("Billing repository not found.")
	} else {
		l.Error("Failed to delete billing repository.", err.Error())
		return http.StatusInternalServerError, errors.New("Failed to delete billing repository.")
	}

}

func getBillRepository(r *http.Request, a routes.Arguments) (int, interface{}) {
	aa := a[aws.AwsAccountSelection].(aws.AwsAccount)
	tx := a[db.Transaction].(*sql.Tx)
	if brs, err := GetBillRepositoryWithPendingForAwsAccount(tx, aa.Id); err != nil {
		l := jsonlog.LoggerFromContextOrDefault(r.Context())
		l.Error("Failed to get aws account's bill repositories.", map[string]interface{}{
			"user":       a[users.AuthenticatedUser].(users.User),
			"awsAccount": aa,
			"error":      err.Error(),
		})
		return http.StatusInternalServerError, errors.New("Failed to retrieve bill repositories.")
	} else {
		var brwss []BillRepositoryWithStatus
		for _, br := range brs {
			brws, _ := WrapBillRepositoriesWithPendingWithStatus(tx, br)
			brwss = append(brwss, brws)
		}
		return http.StatusOK, brwss
	}
}
