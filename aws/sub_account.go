//   Copyright 2018 MSolution.IO
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

	"github.com/trackit/trackit-server/config"
	"github.com/trackit/trackit-server/models"
)

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
		})
	}
	return subAccounts, err
}

// PutSubAccounts gets AWS sub accounts of an aws accounts and puts it in DB if they don't already exists
func PutSubAccounts(ctx context.Context, account AwsAccount, tx *sql.Tx) error {
	subAccounts, err := getAwsSubAccounts(account)
	if err != nil {
		return err
	}
	alreadyAccounts, err := models.AwsAccountsByUserID(tx, account.UserId)
	if err != nil {
		return err
	}
	for _, sub := range subAccounts {
		already := false
		for _, old := range alreadyAccounts {
			if old.AwsIdentity == sub.AwsIdentity {
				already = true
				break
			}
		}
		if !already {
			sub.CreateAwsAccount(ctx, tx)
		}
	}
	return nil
}
