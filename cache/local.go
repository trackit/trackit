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

	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/routes"
	"github.com/trackit/trackit/users"
)

// formatKey is unique depending on user's AWS' identities (personal + shared accounts)
// or identities passed in arguments and route's data (route's name and arguments)
func formatKey(rdCache *redisCache) {
	rdCache.key = fmt.Sprintf("%x-%x-", md5.Sum([]byte(rdCache.route)), md5.Sum([]byte(rdCache.args)))
	for _, val := range rdCache.awsAccount {
		rdCache.key = fmt.Sprintf("%v%v-", rdCache.key, val)
	}
}

// parseRouteFromUrl store the route name (with the format: "/route") and,
// if there is any, argument in the structure redisCache.
func parseRouteFromUrl(url string, rc *redisCache) {
	// The string passed as parameter as the URL is formatted like this: /route?params=value&params2=value2
	idx := strings.IndexByte(url, '?')
	if idx == -1 {
		rc.route = url
	} else {
		// The main purpose is to sort arguments by alphabetical order.
		rc.route = url[:idx]                          // Cut the first part from the URL which is the route
		sortedArgs := strings.Split(url[idx+1:], "&") // Extract, as a strings array, every arguments with their respective values separate by a '&'
		sort.Strings(sortedArgs)                      // Sort the strings by alphabetical order, it's important for the key
		rc.args = strings.Join(sortedArgs, "&")       // Convert from the an array of strings to a single strings and join every part by a '&' like initially
	}
}

// getAwsIdentityFromSharedAcc retrieves all shared accounts based on the user's ID.
// Then the AWS identity, from each respective account, is added to the list of
// AWS identities concerned by the cache.
func getAwsIdentityFromSharedAcc(user users.User, identities *[]string, context *sql.Tx, logger jsonlog.Logger) error {
	sharedAcc, err := models.SharedAccountsByUserID(db.Db, user.Id)
	if err != nil {
		logger.Error("Unable to retrieve AWS' shared accounts by user id.", map[string]interface{}{
			"error":  err.Error(),
			"userId": user.Id,
		})
		return err
	}
	for _, sharedAccContent := range sharedAcc {
		localAcc, localErr := models.AwsAccountByID(context, sharedAccContent.AccountID)
		if localErr != nil {
			logger.Error("Unable to retrieve AWS' account by shared user id.", map[string]interface{}{
				"error":     localErr.Error(),
				"accountID": sharedAccContent.AccountID,
			})
			return localErr
		}
		*identities = append(*identities, localAcc.AwsIdentity)
	}
	return nil
}

// Initialize cache information by getting a list of all AWS identities and
// retrieving different information from the URL. The user's key is also formatted
// depending of the previous information.
func initialiseCacheInfos(url string, args routes.Arguments, logger jsonlog.Logger) (rtn redisCache, err error) {
	var allAcc []string
	parseRouteFromUrl(url, &rtn)
	if args[routes.AwsAccountsOptionalQueryArg] != nil {
		allAcc = args[routes.AwsAccountsOptionalQueryArg].([]string)
	} else {
		tx := args[db.Transaction].(*sql.Tx)
		user := args[users.AuthenticatedUser].(users.User)
		awsAccs, awsAccsErr := models.AwsAccountsByUserID(tx, user.Id)
		if awsAccsErr != nil {
			logger.Error("Unable to retrieve AWS' accounts by user id.", map[string]interface{}{
				"error":  awsAccsErr.Error(),
				"userId": user.Id,
			})
			err = awsAccsErr
			return
		}
		for _, userAccContent := range awsAccs {
			allAcc = append(allAcc, userAccContent.AwsIdentity)
		}
		if err = getAwsIdentityFromSharedAcc(user, &allAcc, tx, logger); err != nil {
			return
		}
	}
	rtn.awsAccount = append(rtn.awsAccount, allAcc...)
	sort.Strings(rtn.awsAccount)
	formatKey(&rtn)
	return
}
