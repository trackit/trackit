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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/config"
	"github.com/trackit/trackit2/models"
	"github.com/trackit/trackit2/users"
)

// AwsAccount represents a client's AWS account.
type AwsAccount struct {
	Id       int    `json:"id"`
	UserId   int    `json:"-"`
	Pretty   string `json:"pretty"`
	RoleArn  string `json:"roleArn"`
	External string `json:"-"`
}

const (
	// assumeRoleDuration is the duration in seconds assumed-role
	// credentials are requested to last. 3600 seconds is the maximum valid
	// value.
	assumeRoleDuration = 3600
)

var (
	// Session is an AWS API session.
	Session client.ConfigProvider
	// stsService gives access to the AWS STS API.
	stsService *sts.STS
	// accountId is the AWS account ID for the credentials provided to the
	// server at startup through the standard methods for the AWS SDK.
	accountId string
)

func init() {
	Session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion),
	}))
	stsService = sts.New(Session)
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
	dbAwsAccounts, err := models.AwsAccountsByUserID(tx, u.Id)
	if err == nil {
		awsAccounts := make([]AwsAccount, len(dbAwsAccounts))
		for i := range dbAwsAccounts {
			awsAccounts[i] = awsAccountFromDbAwsAccount(*dbAwsAccounts[i])
		}
		return awsAccounts, nil
	}
	return nil, err
}

// GetAwsAccountWithIdFromUser returns a user's AWS accounts if it belongs to
// the user.
func GetAwsAccountWithIdFromUser(u users.User, aaid int, tx *sql.Tx) (AwsAccount, error) {
	var aa AwsAccount
	if dbaa, err := models.AwsAccountByID(tx, aaid); err != nil {
		return aa, err
	} else if dbaa.UserID != u.Id {
		return aa, errors.New("aws account does not belong to the user")
	} else {
		aa = awsAccountFromDbAwsAccount(*dbaa)
		return aa, nil
	}
}

// CreateAwsAccount registers a new AWS account for a user. It does no error
// checking: the caller should check themselves that the role ARN exists and is
// correctly configured.
func (a *AwsAccount) CreateAwsAccount(ctx context.Context, db models.XODB) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAwsAccount := models.AwsAccount{
		UserID:  a.UserId,
		RoleArn: a.RoleArn,
		Pretty:  a.Pretty,
		External: sql.NullString{
			Valid:  a.External != "",
			String: a.External,
		},
	}
	err := dbAwsAccount.Insert(db)
	if err == nil {
		a.Id = dbAwsAccount.ID
	} else {
		logger.Error("Failed to insert AWS account in database.", nil)
	}
	return err
}

// awsAccountFromDbAwsAccount constructs an aws.AwsAccount from a
// models.AwsAccount. The distinction exists to decouple database access from
// the logic of the server.
func awsAccountFromDbAwsAccount(dbAwsAccount models.AwsAccount) AwsAccount {
	return AwsAccount{
		Id:       dbAwsAccount.ID,
		UserId:   dbAwsAccount.UserID,
		Pretty:   dbAwsAccount.Pretty,
		RoleArn:  dbAwsAccount.RoleArn,
		External: dbAwsAccount.External.String,
	}
}
