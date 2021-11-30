package shared_account

import (
	"context"
	"database/sql"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/models"
	"github.com/trackit/trackit/users"
)

const (
	adminLevel    = 0
	standardLevel = 1
	readLevel     = 2
)

// safetyCheckByAccountId checks by AccountId if the user have a high enough
// permission level to perform an action on a shared account
func safetyCheckByAccountId(ctx context.Context, tx *sql.Tx, AccountId int, user users.User) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAwsAccount, err := models.AwsAccountByID(tx, AccountId)
	if err == sql.ErrNoRows {
		return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseItemNotFound, "This AWS Account does not exist"})
	} else if err != nil {
		logger.Error("Error while retrieving AWS account from DB", err)
		return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseGenericError, err.Error()})
	}
	if dbAwsAccount.UserID == user.Id {
		return true, nil
	}
	dbSharedAccount, err := models.SharedAccountByAccountID(tx, AccountId)
	if err == nil {
		for _, key := range dbSharedAccount {
			if key.UserID == user.Id && (key.UserPermission == adminLevel || key.UserPermission == standardLevel) {
				return true, nil
			}
		}
	}
	if err == sql.ErrNoRows {
		return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseItemNotFound, "This Shared Account does not exist"})
	}
	logger.Error("Error while retrieving shared account by account ID from DB", err)
	return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseGenericError, "Unable to ensure user have enough rights to do this action"})
}

// checkLevel checks if the current user permission level is high enough to perform an action
func checkLevel(PermissionLevelToCheck int, currentUserPermissionLevel int) bool {
	if currentUserPermissionLevel == adminLevel {
		return true
	} else if currentUserPermissionLevel == standardLevel {
		if currentUserPermissionLevel <= PermissionLevelToCheck {
			return true
		}
	}
	return false
}

// safetyCheckByAccountIdAndPermissionLevel checks by AccountId if the user have a high enough
// permission level to perform an action on a shared account. It also compares current user Permission Level
// to the permissionLevel of the viewer account.
func safetyCheckByAccountIdAndPermissionLevel(ctx context.Context, tx *sql.Tx, AccountId int, body InviteUserRequest, user users.User) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbAwsAccount, err := models.AwsAccountByID(tx, AccountId)
	if err == sql.ErrNoRows {
		return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseItemNotFound, "This AWS Account does not exist"})
	} else if err != nil {
		logger.Error("Error while retrieving AWS account from DB", err)
		return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseGenericError, err.Error()})
	}
	if dbAwsAccount.UserID == user.Id {
		return true, nil
	}
	dbSharedAccount, err := models.SharedAccountByAccountID(tx, AccountId)
	if err == nil {
		for _, key := range dbSharedAccount {
			if key.UserID == user.Id && checkLevel(body.PermissionLevel, key.UserPermission) {
				return true, nil
			}
		}
	}
	if err == sql.ErrNoRows {
		return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseItemNotFound, "This Shared Account does not exist"})
	}
	logger.Error("Error while retrieving shared account by account ID from DB", err)
	return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseGenericError, "Unable to ensure user have enough rights to do this action"})
}

// safetyCheckByShareId checks by ShareId if the user have a high enough
// permission level to perform an action on a shared account
func safetyCheckByShareId(ctx context.Context, tx *sql.Tx, shareId int, user users.User) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbShareAccount, err := models.SharedAccountByID(tx, shareId)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		logger.Error("Error while retrieving Shared Accounts", err)
		return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseGenericError, err.Error()})
	}
	dbAwsAccount, err := models.AwsAccountByID(tx, dbShareAccount.AccountID)
	if err != nil {
		logger.Error("Error while retrieving AWS account", err)
		return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseGenericError, err.Error()})
	}
	if dbAwsAccount.UserID == user.Id {
		return true, nil
	}
	dbShareAccountByAccountId, err := models.SharedAccountByAccountID(tx, dbShareAccount.AccountID)
	if err == nil {
		for _, key := range dbShareAccountByAccountId {
			if key.UserID == user.Id && checkLevel(dbShareAccount.UserPermission, key.UserPermission) {
				return true, nil
			}
		}
	}
	logger.Error("Error while retrieving shared account by account ID from DB", err)
	return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseGenericError, err.Error()})
}

// safetyCheckByShareIdAndPermissionLevel checks by ShareId if the user have a high enough
// permission level to perform an action on a shared account. It also checks if permission Level  of the current user
// is higher than the one that the users wants to set.
func safetyCheckByShareIdAndPermissionLevel(ctx context.Context, tx *sql.Tx, shareId int, newPermissionLevel int, user users.User) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	dbShareAccount, err := models.SharedAccountByID(tx, shareId)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		logger.Error("Error while retrieving Shared Accounts from DB", err)
		return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseGenericError, err.Error()})
	}
	dbAwsAccount, err := models.AwsAccountByID(tx, dbShareAccount.AccountID)
	if err != nil {
		logger.Error("Error while retrieving AWS account", err)
		return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseGenericError, err.Error()})
	}
	if dbAwsAccount.UserID == user.Id {
		return true, nil
	}
	dbShareAccountByAccountId, err := models.SharedAccountByAccountID(tx, dbShareAccount.AccountID)
	if err == nil {
		for _, key := range dbShareAccountByAccountId {
			if key.UserID == user.Id && checkLevel(newPermissionLevel, key.UserPermission) {
				if key.UserPermission == standardLevel && dbShareAccount.UserPermission == adminLevel {
					return false, nil
				}
				return true, nil
			}
		}
	}
	logger.Error("Error while retrieving shared account by account ID from DB", err)
	return false, errors.GetErrorMessage(ctx, &errors.DatabaseError{errors.DatabaseGenericError, err.Error()})
}

// checkPermissionLevel checks user permission level
func checkPermissionLevel(permissionLevel int) bool {
	if permissionLevel == adminLevel {
		return true
	} else if permissionLevel == standardLevel {
		return true
	} else if permissionLevel == readLevel {
		return true
	} else {
		return false
	}
}
