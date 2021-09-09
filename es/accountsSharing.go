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

package es

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/users"
)

// AccountsAndIndexes stores the accounts and indexes
type AccountsAndIndexes struct {
	Accounts []string
	Indexes  []string
}

// isAccountDuplicate returns true if the account already exists in the list of accounts
// returns false otherwise
func (ai *AccountsAndIndexes) isAccountDuplicate(account string) bool {
	for _, entry := range ai.Accounts {
		if entry == account {
			return true
		}
	}
	return false
}

// addAccount adds a new account in the AccountsAndIndexes if it is not already
// in the list of accounts
func (ai *AccountsAndIndexes) addAccount(account string) {
	if !ai.isAccountDuplicate(account) {
		ai.Accounts = append(ai.Accounts, account)
	}
}

// addIndex adds a new index in the AccountsAndIndexes if it is not already
// in the list of indexes
func (ai *AccountsAndIndexes) addIndex(index string) {
	for _, entry := range ai.Indexes {
		if entry == index {
			return
		}
	}
	ai.Indexes = append(ai.Indexes, index)
}

// getAllAccountsAndIndexes returns an AccountsAndIndexes struct, a status code and an error
// The AccountsAndIndexes struct will contain all the accounts available to the user
// with their indexes without duplicates
// If an account is shared with a user that already have the same account, the shared
// account and index will be skipped
func getAllAccountsAndIndexes(user users.User, tx *sql.Tx, indexPrefix string) (AccountsAndIndexes, int, error) {
	accountsAndIndexes := AccountsAndIndexes{}
	// Retrieve the user accounts and shared accounts
	userAccounts, err := models.AwsAccountsByUserID(tx, user.Id)
	if err != nil {
		return accountsAndIndexes, http.StatusInternalServerError, fmt.Errorf("Unable to retrieve the list of accounts for current user: %s", err.Error())
	}
	sharedAccounts, err := models.SharedAccountsWithRoleByUserID(tx, user.Id)
	if err != nil {
		return accountsAndIndexes, http.StatusInternalServerError, fmt.Errorf("Unable to retrieve the list of shared accounts for current user: %s", err.Error())
	}
	// Add all the user accounts
	for _, userAccount := range userAccounts {
		accountsAndIndexes.addAccount(userAccount.AwsIdentity)
		accountsAndIndexes.addIndex(IndexNameForUserId(userAccount.UserID, indexPrefix))
	}
	// Add all the non duplicate shared accounts
	for _, sharedAccount := range sharedAccounts {
		// Do not add the account if the user already own the same account
		if !accountsAndIndexes.isAccountDuplicate(sharedAccount.AwsIdentity) {
			accountsAndIndexes.addAccount(sharedAccount.AwsIdentity)
			accountsAndIndexes.addIndex(IndexNameForUserId(sharedAccount.OwnerID, indexPrefix))
		}
	}
	// If no indexes where found, return an error to prevent giving access to all indexes
	if len(accountsAndIndexes.Indexes) == 0 {
		return accountsAndIndexes, http.StatusBadRequest, fmt.Errorf("No aws account found")
	}
	return accountsAndIndexes, http.StatusOK, nil
}

// GetAccountsAndIndexes returns an AccountsAndIndexes struct, a status code and an error
// if the accountList parameter is empty the function will call getAllAccountsAndIndexes
// if the accountList parameter is not empty the function will validate the accounts and
// find their indexes
func GetAccountsAndIndexes(accountList []string, user users.User, tx *sql.Tx, indexPrefix string) (AccountsAndIndexes, int, error) {
	if len(accountList) == 0 {
		return getAllAccountsAndIndexes(user, tx, indexPrefix)
	}
	accountsAndIndexes := AccountsAndIndexes{}
	if err := aws.ValidateAwsAccounts(accountList); err != nil {
		return accountsAndIndexes, http.StatusBadRequest, err
	}
	// Retrieve the user's accounts and shared accounts
	userAccounts, err := models.AwsAccountsByUserID(tx, user.Id)
	if err != nil {
		return accountsAndIndexes, http.StatusInternalServerError, fmt.Errorf("Unable to retrieve the list of accounts for current user: %s", err.Error())
	}
	sharedAccounts, err := models.SharedAccountsWithRoleByUserID(tx, user.Id)
	if err != nil {
		return accountsAndIndexes, http.StatusInternalServerError, fmt.Errorf("Unable to retrieve the list of shared accounts for current user: %s", err.Error())
	}
	// Match the accountList parameter with the user's accounts and shared accounts
	for _, account := range accountList {
		found_match := false
		// Try to match in priority with the user's accounts
		for _, userAccount := range userAccounts {
			if userAccount.AwsIdentity == account {
				found_match = true
				accountsAndIndexes.addAccount(userAccount.AwsIdentity)
				accountsAndIndexes.addIndex(IndexNameForUserId(userAccount.UserID, indexPrefix))
			}
		}
		// If no match is found in the user's accounts, try in the shared accounts
		if !found_match {
			for _, sharedAccount := range sharedAccounts {
				if sharedAccount.AwsIdentity == account {
					found_match = true
					if !accountsAndIndexes.isAccountDuplicate(sharedAccount.AwsIdentity) {
						accountsAndIndexes.addAccount(sharedAccount.AwsIdentity)
						accountsAndIndexes.addIndex(IndexNameForUserId(sharedAccount.OwnerID, indexPrefix))
					}
				}
			}
		}
		if !found_match {
			return accountsAndIndexes, http.StatusBadRequest, fmt.Errorf("Unable to access account %s", account)
		}
	}
	return accountsAndIndexes, http.StatusOK, nil
}
