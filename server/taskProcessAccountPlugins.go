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

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/cache"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/accountPlugins"
	"github.com/trackit/trackit/models"
	core "github.com/trackit/trackit/plugins/account/core"
	"github.com/trackit/trackit/users"
)

// taskProcessAccountPlugins is the entry point for account plugins processing
func taskProcessAccountPlugins(ctx context.Context) error {
	args := paramsFromContextOrArgs(ctx)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'process-account-plugins'.", map[string]interface{}{
		"args": args,
	})
	if len(args) != 1 {
		return errors.New("taskProcessAccountPlugins requires an integer argument")
	} else if aaId, err := strconv.Atoi(args[0]); err != nil {
		return err
	} else {
		return preparePluginsProcessingForAccount(ctx, aaId)
	}
}

// preparePluginsProcessingForAccount retrieves all the information needed to
// run the plugins for a given account
func preparePluginsProcessingForAccount(ctx context.Context, aaId int) (err error) {
	var tx *sql.Tx
	var aa aws.AwsAccount
	var user users.User
	var updateId int64
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	defer utilsUsualTxFinalize(&tx, &err, &logger, "process-account-plugins")

	var trackitUser *models.User // We can't use := because then there would be a new err which would shadow the returned value
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else if aa, err = aws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else if trackitUser, err = models.UserByID(db.Db, aa.UserId); err != nil || trackitUser.AccountType != "trackit" {
		if err == nil {
			logger.Info("Task 'ProcessAccountPlugins' has been skipped because the user has the wrong account type.", map[string]interface{}{
				"userAccountType": trackitUser.AccountType,
				"requiredAccount": "trackit",
			})
		}
	} else if user, err = users.GetUserWithId(tx, aa.UserId); err != nil {
	} else if updateId, err = registerAccountPluginsProcessing(db.Db, aa); err != nil {
	} else {
		runPluginsForAccount(ctx, user, aa)
		updateAccountPluginsCompletion(ctx, aaId, db.Db, updateId, nil)
	}
	if err != nil {
		updateAccountPluginsCompletion(ctx, aaId, db.Db, updateId, err)
		logger.Error("Failed to process account plugins.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
	}
	var affectedRoutes = []string{
		"/plugins/results",
	}
	err = cache.RemoveMatchingCache(affectedRoutes, []string{aa.AwsIdentity}, logger)
	return
}

// runPluginsForAccount runs all the registered plugins for an account
func runPluginsForAccount(ctx context.Context, user users.User, aa aws.AwsAccount) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	for _, plugin := range core.RegisteredAccountPlugins {
		if !plugin.BillingDataOnly && strings.TrimSpace(aa.RoleArn) == "" {
			continue
		}
		accountId := aa.AwsIdentity
		pluginResultES := accountPlugins.PluginResultES{
			Account:    accountId,
			ReportDate: time.Now().UTC(),
			PluginName: plugin.Name,
			Category:   plugin.Category,
			Label:      plugin.Label,
		}
		params := core.PluginParams{
			Context:    ctx,
			User:       user,
			AwsAccount: aa,
			AccountId:  accountId,
			ESClient:   es.Client,
		}
		if !plugin.BillingDataOnly {
			creds, err := aws.GetTemporaryCredentials(aa, fmt.Sprintf("trackit-%s-plugin", plugin.Name))
			if err != nil {
				logger.Error("Error when getting temporary credentials", err.Error())
				pluginResultES.Error = fmt.Sprintf("Error when getting temporary credentials: %s", err.Error())
			} else {
				params.AccountCredentials = creds
			}
		}
		if pluginResultES.Error == "" {
			res := plugin.Func(params)
			pluginResultES.Result = res.Result
			pluginResultES.Status = res.Status
			pluginResultES.Details = res.Details
			pluginResultES.Error = res.Error
			pluginResultES.Checked = res.Checked
			pluginResultES.Passed = res.Passed
		}
		core.IngestPluginResult(ctx, aa, pluginResultES)
	}
}

func registerAccountPluginsProcessing(db *sql.DB, aa aws.AwsAccount) (int64, error) {
	const sqlstr = `INSERT INTO aws_account_plugins_job(
	aws_account_id,
	worker_id
	) VALUES (?, ?)`
	res, err := db.Exec(sqlstr, aa.Id, backendId)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func updateAccountPluginsCompletion(ctx context.Context, aaId int, db *sql.DB, updateId int64, jobErr error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	doErr := func(errorText string, err error) {
		logger.Error(errorText, map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
			"updateId":     updateId,
		})
	}
	if uErr := updateNextUpdateAccountPlugins(db, aaId); uErr != nil {
		doErr("Failed to update the account next plugin update date", uErr)
	}
	if rErr := registerAccountPluginsCompletion(db, updateId, jobErr); rErr != nil {
		doErr("Failed to register account plugins completion", rErr)
	}
}

func updateNextUpdateAccountPlugins(db *sql.DB, aaId int) error {
	const sqlstr = `UPDATE aws_account SET
	next_update_plugins=?
	WHERE id=?`
	_, err := db.Exec(sqlstr, time.Now().AddDate(0, 0, 1), aaId)
	return err
}

func registerAccountPluginsCompletion(db *sql.DB, updateId int64, jobErr error) error {
	const sqlstr = `UPDATE aws_account_plugins_job SET
	completed=?,
	jobError=?
	WHERE id=?`
	jobError := ""
	if jobErr != nil {
		jobError = jobErr.Error()
	}
	_, err := db.Exec(sqlstr, time.Now(), jobError, updateId)
	return err
}
