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
	"database/sql"
	"errors"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

// uintSliceContainsUint returns true if a slice of uint contains id.
func uintSliceContainsUint(slice []uint, id uint) bool {
	for _, idd := range slice {
		if idd == id {
			return true
		}
	}
	return false
}

// stripAwsAccounts strips the slice of AwsAccount with filters in the
// arguments map with QueryArgAwsAccounts as the key.
func stripAwsAccounts(awsAccounts []AwsAccount, a routes.Arguments) []AwsAccount {
	if filters, ok := a[QueryArgAwsAccounts]; ok {
		for index, awsAccount := range awsAccounts {
			if !uintSliceContainsUint(filters.([]uint), uint(awsAccount.Id)) {
				length := len(awsAccounts)
				awsAccounts[index] = awsAccounts[length-1]
				awsAccounts = awsAccounts[:length-1]
			}
		}
	}
	return awsAccounts
}

// getAwsAccount is a route handler which returns the caller's list of
// AwsAccounts. It handles QueryArgAwsAccounts thanks to stripAwsAccounts.
func getAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	u := a[users.AuthenticatedUser].(users.User)
	tx := a[db.Transaction].(*sql.Tx)
	l := jsonlog.LoggerFromContextOrDefault(r.Context())
	awsAccounts, err := GetAwsAccountsFromUser(u, tx)
	if err == nil {
		return 200, stripAwsAccounts(awsAccounts, a)
	} else {
		l.Error("Failed to get user's AWS accounts.", err.Error())
		return 500, errors.New("Failed to retrieve AWS accounts.")
	}
}
