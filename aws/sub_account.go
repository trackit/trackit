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
	"database/sql"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"

	"github.com/trackit/trackit-server/config"
	"github.com/trackit/trackit-server/models"
)

type SubAccount struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// GetSubAccountByAwsAccountId gets SubAccounts of an aws account from the DB
func GetSubAccountsByAwsAccountId(aaId int, tx *sql.Tx) ([]SubAccount, error) {
	subAccounts, err := models.AwsSubAccountsByAwsAccountID(tx, aaId)
	if err != nil {
		return nil, err
	}
	response := make([]SubAccount, 0)
	for _, subAccount := range subAccounts {
		response = append(response, SubAccount{
			Id:   subAccount.AwsID,
			Name: subAccount.Name,
		})
	}
	return response, nil
}

// addNewSubAccounts adds new sub accounts for an aws account in DB
func addNewSubAccounts(oldSubAccounts []*models.AwsSubAccount, newSubAccounts []SubAccount, aaId int, awsId string, tx *sql.Tx) error {
	for _, newSubAccount := range newSubAccounts {
		already := false
		for _, old := range oldSubAccounts {
			if old.AwsID == newSubAccount.Id && old.Name == newSubAccount.Name {
				already = true
				break
			}
		}
		if !already && newSubAccount.Id != awsId {
			asa := &models.AwsSubAccount{
				AwsAccountID: aaId,
				AwsID:        newSubAccount.Id,
				Name:         newSubAccount.Name,
			}
			err := asa.Insert(tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// rmOldSubAccounts removes sub accounts that aren't sub accounts for an aws account in DB
func rmOldSubAccounts(oldSubAccounts []*models.AwsSubAccount, newSubAccounts []SubAccount, aaId int, tx *sql.Tx) error {
	for _, old := range oldSubAccounts {
		keep := false
		for _, newSubAccount := range newSubAccounts {
			if old.AwsID == newSubAccount.Id && old.Name == newSubAccount.Name {
				keep = true
				break
			}
		}
		if !keep {
			err := old.Insert(tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetAwsSubAccounts gets the list of sub accounts in AWS for an aws account
func GetAwsSubAccounts(aa AwsAccount) ([]SubAccount, error) {
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
	subAccounts := make([]SubAccount, 0)
	for _, account := range res.Accounts {
		subAccounts = append(subAccounts, SubAccount{
			Id:   aws.StringValue(account.Id),
			Name: aws.StringValue(account.Name),
		})
	}
	return subAccounts, err
}

// UpdateSubAccounts updates sub accounts of an aws account thanks to AWS
func UpdateSubAccounts(aa AwsAccount, tx *sql.Tx) error {
	split := strings.Split(aa.RoleArn, ":")
	awsId := ""
	if len(split) >= 5 {
		awsId = split[4]
	}
	oldSubAccounts, err := models.AwsSubAccountsByAwsAccountID(tx, aa.Id)
	if err != nil {
		return err
	}
	newSubAccounts, err := GetAwsSubAccounts(aa)
	if err != nil {
		return err
	}
	err = addNewSubAccounts(oldSubAccounts, newSubAccounts, aa.Id, awsId, tx)
	if err != nil {
		return err
	}
	err = rmOldSubAccounts(oldSubAccounts, newSubAccounts, aa.Id, tx)
	if err != nil {
		return err
	}
	return nil
}
