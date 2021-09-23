//   Copyright 2019 MSolution.IO
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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/models"
)

func updateExistingAccount(ctx context.Context, aa AwsAccount, subs []AwsAccount, tx *sql.Tx) error {
	exists, err := models.AwsAccountByUserID(tx, aa.UserId)
	if err != nil {
		return err
	}
	for _, exist := range exists {
		if exist.ParentID.Valid {
			continue
		}
		for _, sub := range subs {
			if exist.AwsIdentity == sub.AwsIdentity && exist.AwsIdentity != aa.AwsIdentity {
				exist.ParentID.Valid = true
				exist.ParentID.Int64 = int64(aa.Id)
				if updateErr := exist.Update(tx); updateErr != nil {
					if err == nil {
						err = updateErr
					}
					jsonlog.LoggerFromContextOrDefault(ctx).Error("Failure to update AWS account in database", map[string]interface{}{
						"error": updateErr.Error(),
					})
				}
			}
		}
	}
	return err
}

// getAwsSubAccounts gets the list of sub accounts in AWS for an aws account
func getAwsSubAccounts(aa AwsAccount) ([]AwsAccount, error) {
	creds, err := GetTemporaryCredentials(aa, "GetAwsSubAccounts")
	if err != nil {
		return nil, err
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	orga := organizations.New(sess)
	res, err := orga.ListAccounts(nil)
	if err != nil {
		return nil, err
	}
	subAccounts := make([]AwsAccount, 0)
	for _, account := range res.Accounts {
		subAccounts = append(subAccounts, AwsAccount{
			UserId:      aa.UserId,
			Pretty:      aws.StringValue(account.Name),
			RoleArn:     "",
			External:    "",
			Payer:       false,
			AwsIdentity: aws.StringValue(account.Id),
			ParentId:    sql.NullInt64{int64(aa.Id), true},
		})
	}
	return subAccounts, err
}

// PutSubAccounts gets AWS sub accounts of an aws accounts and puts it in DB if they don't already exists
func PutSubAccounts(ctx context.Context, account AwsAccount, tx *sql.Tx) error {
	identity, err := account.GetAwsAccountIdentity()
	if err != nil {
		return err
	}
	account.AwsIdentity = identity
	subAccounts, err := getAwsSubAccounts(account)
	if err != nil {
		return err
	}
	alreadyAccounts, err := models.AwsAccountByUserID(tx, account.UserId)
	if err != nil {
		return err
	}
SubAccountsLoop:
	for _, sub := range subAccounts {
		for _, old := range alreadyAccounts {
			if old.AwsIdentity == sub.AwsIdentity {
				continue SubAccountsLoop
			}
		}
		if createError := sub.CreateAwsAccount(ctx, tx); createError != nil {
			if err == nil {
				err = createError
			}
			jsonlog.LoggerFromContextOrDefault(ctx).Error("Failure to register AWS sub-account", map[string]interface{}{
				"error": createError.Error(),
			})
		}
	}
	if err != nil {
		return err
	}
	return updateExistingAccount(ctx, account, subAccounts, tx)
}
