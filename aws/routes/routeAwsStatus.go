//   Copyright 2021 MSolution.IO
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
package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

func getAwsAccountsStatus(r *http.Request, a routes.Arguments) (int, interface{}) {
	var awsAccounts []aws.AwsAccount
	var awsAccountsWithBillRepositories []s3.AwsAccountWithBillRepositoriesWithPending
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	awsAccounts, err := aws.GetAwsAccountsFromUser(u, tx)
	if err != nil {
		l.Error("failed to get user's AWS accounts", err.Error())
		return http.StatusInternalServerError, errors.New("failed to retrieve AWS accounts")
	}
	awsAccountsWithBillRepositories, err = s3.WrapAwsAccountsWithBillRepositories(awsAccounts, tx)
	if err != nil {
		l.Error("failed to get AWS accounts' bill repositories", err.Error())
		return http.StatusInternalServerError, errors.New("failed to retrieve bill repositories")
	}
	// Code unneeded for now considering all it does is create an array and fill it with all the id data from the AWS stuff before proceeding to do nothing with it
	/*
		billRepositoriesIds := make([]int, 0)
		for _, awsAccount := range awsAccountsWithBillRepositories {
			for _, billRepository := range awsAccount.BillRepositories {
				billRepositoriesIds = append(billRepositoriesIds, billRepository.Id)
			}
		}
	*/
	result := s3.WrapAwsAccountsWithBillRepositoriesWithPendingWithStatus(awsAccountsWithBillRepositories, tx)
	return http.StatusOK, result
}
