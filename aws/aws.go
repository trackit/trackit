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

package aws

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/awsSession"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/users"
)

// AwsAccount represents a client's AWS account.
type AwsAccount struct {
	Id             int           `json:"id"`
	UserId         int           `json:"-"`
	Pretty         string        `json:"pretty"`
	RoleArn        string        `json:"roleArn"`
	External       string        `json:"-"`
	Payer          bool          `json:"payer"`
	AccountOwner   bool          `json:"accountOwner"`
	UserPermission int           `json:"permissionLevel"`
	AwsIdentity    string        `json:"awsIdentity"`
	ParentId       sql.NullInt64 `json:"-"`
}

var (
	// stsService gives access to the AWS STS API.
	stsService *sts.STS
	// accountId is the AWS account ID for the credentials provided to the
	// server at startup through the standard methods for the AWS SDK.
	accountId string
)

func init() {
	stsService = sts.New(awsSession.Session)
	accountId = initAccountId(stsService)
}

// initAccountId uses the AWS STS API's GetCallerIdentity method to discover
// the account ID for the AWS account the server uses.
func initAccountId(s *sts.STS) string {
	var input sts.GetCallerIdentityInput
	output, err := s.GetCallerIdentity(&input)
	if err != nil {
		log.Fatalf("Failed to get AWS account ID: '%s'.", err.Error())
	}
	return *output.Account
}

// AccountId returns the server's AWS account ID.
func AccountId() string { return accountId }

// GetAwsAccountFromUser returns a slice of all AWS accounts configured by a
// given user.
func GetAwsAccountsFromUser(u users.User, tx *sql.Tx) ([]AwsAccount, error) {
	var res []AwsAccount
	dbAwsAccounts, err := models.AwsAccountsByUserID(tx, u.Id)
	if err != nil {
		return nil, err
	}
	dbShareAccounts, err := models.SharedAccountsByUserID(tx, u.Id)
	if err != nil {
		return nil, err
	}
	for _, key := range dbAwsAccounts {
		res = append(res, AwsAccount{
			key.ID,
			key.UserID,
			key.Pretty,
			key.RoleArn,
			key.External,
			key.Payer,
			true,
			0,
			key.AwsIdentity,
			key.ParentID})
	}
	for _, key := range dbShareAccounts {
		dbAwsAccountById, err := models.AwsAccountByID(tx, key.AccountID)
		if err != nil {
			return nil, err
		}
		res = append(res, AwsAccount{
			dbAwsAccountById.ID,
			dbAwsAccountById.UserID,
			dbAwsAccountById.Pretty,
			dbAwsAccountById.RoleArn,
			dbAwsAccountById.External,
			dbAwsAccountById.Payer,
			false,
			key.UserPermission,
			dbAwsAccountById.AwsIdentity,
			dbAwsAccountById.ParentID})
	}
	return res, nil
}

// GetAwsAccountWithId returns an AWS account.
func GetAwsAccountWithId(aaid int, tx *sql.Tx) (AwsAccount, error) {
	var aa AwsAccount
	if dbaa, err := models.AwsAccountByID(tx, aaid); err != nil {
		return aa, err
	} else {
		aa = AwsAccountFromDbAwsAccount(*dbaa)
		return aa, nil
	}
}

// GetAwsAccountWithIdFromUser returns a user's AWS accounts if it belongs to
// the user.
func GetAwsAccountWithIdFromUser(u users.User, aaid int, tx *sql.Tx) (AwsAccount, error) {
	var aaz AwsAccount
	if aa, err := GetAwsAccountWithId(aaid, tx); err != nil {
		return aaz, err
	} else if aa.UserId == u.Id {
		return aa, nil
	} else {
		return aaz, errors.New("aws account does not belong to the user")
	}
}

// CreateAwsAccount registers a new AWS account for a user. It does no error
// checking: the caller should check themselves that the role ARN exists and is
// correctly configured.
func (a *AwsAccount) CreateAwsAccount(ctx context.Context, db models.XODB) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identity, err := a.GetAwsAccountIdentity()
	if err != nil {
		logger.Error("Failed to get AWS Identity.", err.Error())
		return err
	}
	dbAwsAccount := models.AwsAccount{
		UserID:      a.UserId,
		RoleArn:     a.RoleArn,
		Pretty:      a.Pretty,
		External:    a.External,
		Payer:       a.Payer,
		AwsIdentity: identity,
		ParentID:    a.ParentId,
	}
	err = dbAwsAccount.Insert(db)
	if err == nil {
		a.Id = dbAwsAccount.ID
	} else {
		logger.Error("Failed to insert AWS account in database.", nil)
	}
	return err
}

// UpdatePrettyAwsAccount updates an AWS account for a user. It does no error
// checking: the caller should check themselves that the AWS account exists.
// Only the Pretty will be updated.
func (a *AwsAccount) UpdatePrettyAwsAccount(ctx context.Context, tx *sql.Tx) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAwsAccount, err := models.AwsAccountByID(tx, a.Id)
	if err != nil {
		logger.Error("Failed to get AWS account in database.", err.Error())
	} else {
		dbAwsAccount.Pretty = a.Pretty
		dbAwsAccount.Payer = a.Payer
		dbAwsAccount.RoleArn = a.RoleArn
		err := dbAwsAccount.Update(tx)
		if err != nil {
			logger.Error("Failed to update AWS account in database.", err.Error())
		}
	}
	return err
}

// UpdateIdentityAwsAccount updates an AWS account for a user. It does no error
// checking: the caller should check themselves that the AWS account exists.
// Only the identity will be updated.
func (a *AwsAccount) UpdateIdentityAwsAccount(ctx context.Context, tx *sql.Tx) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAwsAccount, err := models.AwsAccountByID(tx, a.Id)
	if err != nil {
		logger.Error("Failed to get AWS account in database.", err.Error())
		return err
	}
	identity, err := a.GetAwsAccountIdentity()
	if err != nil {
		logger.Error("Failed to get AWS identity.", err.Error())
		return err
	}
	dbAwsAccount.AwsIdentity = identity
	err = dbAwsAccount.Update(tx)
	if err != nil {
		logger.Error("Failed to update AWS account in database.", err.Error())
	}
	return err
}

// UpdateRoleAndExternalAwsAccount updates an AWS account for a user. It does no error
// checking: the caller should check themselves that the AWS account exists.
// Only the RoleArn and External will be updated.
func (a *AwsAccount) UpdateRoleAndExternalAwsAccount(ctx context.Context, tx *sql.Tx) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAwsAccount, err := models.AwsAccountByID(tx, a.Id)
	if err != nil {
		logger.Error("Failed to get AWS account in database.", err.Error())
	} else {
		dbAwsAccount.RoleArn = a.RoleArn
		dbAwsAccount.External = a.External
		err := dbAwsAccount.Update(tx)
		if err != nil {
			logger.Error("Failed to update AWS account in database.", err.Error())
		}
	}
	return err
}

// GetAwsAccountIdentity returns the AWS identity of an AWS Account.
func (a *AwsAccount) GetAwsAccountIdentity() (identity string, err error) {
	if strings.TrimSpace(a.RoleArn) == "" {
		return a.AwsIdentity, nil
	} else if reg, err := regexp.Compile("^arn:aws:iam::([0-9]{12}):.*$"); err != nil {
		return "", err
	} else {
		identity = reg.FindStringSubmatch(a.RoleArn)[1]
	}
	return
}

// AwsAccountFromDbAwsAccount constructs an aws.AwsAccount from a
// models.AwsAccount. The distinction exists to decouple database access from
// the logic of the server.
func AwsAccountFromDbAwsAccount(dbAwsAccount models.AwsAccount) AwsAccount {
	return AwsAccount{
		Id:          dbAwsAccount.ID,
		UserId:      dbAwsAccount.UserID,
		Pretty:      dbAwsAccount.Pretty,
		RoleArn:     dbAwsAccount.RoleArn,
		External:    dbAwsAccount.External,
		Payer:       dbAwsAccount.Payer,
		AwsIdentity: dbAwsAccount.AwsIdentity,
		ParentId:    dbAwsAccount.ParentID,
	}
}
