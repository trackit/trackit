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

package cache

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/models"
	"github.com/trackit/trackit-server/routes"
	"github.com/trackit/trackit-server/users"
)

// formatKey is unique depending on user's AWS' identities (personal + shared accounts)
// or identities passed in arguments and route's URL
func formatKey(rdCache *redisCache) {
	rdCache.key = fmt.Sprintf("%x-%x:", md5.Sum([]byte(rdCache.route)), md5.Sum([]byte(rdCache.args)))
	for _, val := range rdCache.awsAccount {
		rdCache.key = fmt.Sprintf("%v%v:", rdCache.key, val)
	}
}

func parseRouteFromUrl(url string, rc *redisCache) {
	idx := strings.IndexByte(url, '?')
	if idx == -1 {
		rc.route = url
	} else {
		rc.route = url[:idx]
		rc.args = url[idx + 1:]
		sortedArgs := strings.Split(rc.args, "&")
		sort.Strings(sortedArgs)
		rc.args = strings.Join(sortedArgs, "&")
	}
}

// Retrieving all shared accounts based on the user's ID.
// Then the AWS identity, from each respective account, is added to the list of
// AWS identities concerned by the cache.
func getAwsIdentityFromSharedAcc(user users.User, identities *[]string, context *sql.Tx, logger jsonlog.Logger) error {
	sharedAcc, err := models.SharedAccountsByUserID(db.Db, user.Id)
	if err != nil {
		logger.Error("Unable to retrieve AWS' shared accounts by user id.", map[string] interface{} {
			"error":  err.Error(),
			"userId": user.Id,
		})
		return err
	}
	for _, sharedAccContent := range sharedAcc {
		localAcc, localErr := models.AwsAccountByID(context, sharedAccContent.AccountID)
		if localErr != nil {
			logger.Error("Unable to retrieve AWS' account by shared user id.", map[string] interface{} {
				"error":     localErr.Error(),
				"accountID": sharedAccContent.AccountID,
			})
			return localErr
		}
		*identities = append(*identities, localAcc.AwsIdentity)
	}
	return nil
}

// Listing of all AWS identities concerned by the cache depending
// if an (or multiple) AWS identity are passed in arguments or not.
func retrieveRouteInfos(url string, args routes.Arguments, logger jsonlog.Logger) (rtn redisCache, err error) {
	var allAcc []string
	parseRouteFromUrl(url, &rtn)
	if args[routes.AwsAccountsOptionalQueryArg] != nil {
		allAcc = args[routes.AwsAccountsOptionalQueryArg].([]string)
	} else {
		tx := args[db.Transaction].(*sql.Tx)
		user := args[users.AuthenticatedUser].(users.User)
		awsAccs, awsAccsErr := models.AwsAccountsByUserID(tx, user.Id)
		if awsAccsErr != nil {
			logger.Error("Unable to retrieve AWS' accounts by user id.", map[string] interface{} {
				"error":  awsAccsErr.Error(),
				"userId": user.Id,
			})
			return rtn, awsAccsErr
		}
		for _, userAccContent := range awsAccs {
			allAcc = append(allAcc, userAccContent.AwsIdentity)
		}
		if err = getAwsIdentityFromSharedAcc(user, &allAcc, tx, logger); err != nil {
			return
		}
	}
	for _, val := range allAcc {
		rtn.awsAccount = append(rtn.awsAccount, val)
	}
	sort.Strings(rtn.awsAccount)
	formatKey(&rtn)
	return
}
