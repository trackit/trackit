package shared_account

import (
	"database/sql"
	"errors"
	"context"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/users"
	"github.com/trackit/trackit-server/models"
)

// safetyCheckByAccountId checks by AccountId if the user have a high enough
// permission level to perform an action on a shared account
func safetyCheckByAccountId(ctx context.Context, tx *sql.Tx, AccountId int, user users.User) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAwsAccount, err := models.AwsAccountByID(tx, AccountId)
	if err == sql.ErrNoRows {
		logger.Error("Non existing AWS error", err)
		return false, errors.New("This AWS Account does not exist")
	} else if err != nil {
		logger.Error("Unable to ensure user have enough rights to do this action", err)
		return false, err
	}
	if dbAwsAccount.UserID == user.Id {
		return true, nil
	}
	dbSharedAccount, err := models.SharedAccountsByAccountID(tx, AccountId)
	if err == nil {
		for _, key := range dbSharedAccount {
			if key.UserID == user.Id  && (key.UserPermission == 0 || key.UserPermission == 1){
				return true, nil
			}
		}
	}
	logger.Error("Unable to ensure user have enough rights to do this action", err)
	return false, errors.New("Unable to ensure user have enough rights to do this action")
}

// checkLevel checks if the current user permission level is high enough to perform an action
func checkLevel(sharePermissionLevel int, permissionLevel int) (bool) {
	if permissionLevel == 0 {
		return true
	} else if permissionLevel == 1 {
		if permissionLevel <= sharePermissionLevel {
			return true
		}
	}
	return false
}

// safetyCheckByShareId checks by ShareId if the user have a high enough
// permission level to perform an action on a shared account
func safetyCheckByShareId(ctx context.Context, tx *sql.Tx, shareId int, user users.User) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbShareAccount, err := models.SharedAccountByID(tx, shareId)
	if err == sql.ErrNoRows {
		logger.Error("Non existing Shared access", err)
		return false, nil
	} else if err != nil {
		logger.Error("Error while retrieving Shared Accounts" , err)
		return false, err
	}
	dbAwsAccount, err := models.AwsAccountByID(tx, dbShareAccount.AccountID)
	if dbAwsAccount.UserID == user.Id {
		return true, nil
	}
	dbShareAccountByAccountId, err := models.SharedAccountsByAccountID(tx, dbShareAccount.AccountID)
	if err == nil {
		for _, key := range dbShareAccountByAccountId {
			if key.UserID == user.Id && checkLevel(dbShareAccount.UserPermission, key.UserPermission) {
				return true, nil
			}
		}
	}
	logger.Error("Unable to ensure user have enough rights to do this action", err)
	return false, err
}

// safetyCheckByShareIdAndPermissionLevel checks by ShareId if the user have a high enough
// permission level to perform an action on a shared account. It also checks if permission Level  of the current user
// is higher than the one that the users wants to set.
func safetyCheckByShareIdAndPermissionLevel(ctx context.Context, tx *sql.Tx, shareId int, newPermissionLevel int, user users.User) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbShareAccount, err := models.SharedAccountByID(tx, shareId)
	if err == sql.ErrNoRows {
		logger.Error("Non existing Shared access", err)
		return false, nil
	} else if err != nil {
		logger.Error("Error while retrieving Shared Accounts" , err)
		return false, err
	}
	dbAwsAccount, err := models.AwsAccountByID(tx, dbShareAccount.AccountID)
	if dbAwsAccount.UserID == user.Id {
		return true, nil
	}
	dbShareAccountByAccountId, err := models.SharedAccountsByAccountID(tx, dbShareAccount.AccountID)
	if err == nil {
		for _, key := range dbShareAccountByAccountId {
			if key.UserID == user.Id && checkLevel(newPermissionLevel, key.UserPermission) {
				return true, nil
			}
		}
	}
	logger.Error("Unable to ensure user have enough rights to do this action", err)
	return false, err
}

// checkPermissionLevel checks user permission level
func checkPermissionLevel(permissionLevel int) (bool) {
	if permissionLevel == 0 {
		return true
	} else if permissionLevel == 1 {
		return true
	} else if permissionLevel == 2 {
		return true
	} else {
		return false
	}
}
