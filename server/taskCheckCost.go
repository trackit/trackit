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

package main

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"

	"github.com/trackit/jsonlog"
	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/costs"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/models"
)

// taskCheckCost is the entry point for account cost verification
func taskCheckCost(ctx context.Context) error {
	args := paramsFromContextOrArgs(ctx)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'check-cost'.", map[string]interface{}{
		"args": args,
	})
	if len(args) != 1 {
		return errors.New("taskCheckCost requires an integer argument")
	} else if aaId, err := strconv.Atoi(args[0]); err != nil {
		return err
	} else {
		return prepareCheckCostForAccount(ctx, aaId)
	}
}

// prepareCheckCostForAccount retrieves all the informations needed to
// run a cost check for a given account
func prepareCheckCostForAccount(ctx context.Context, aaId int) (err error) {
	var tx *sql.Tx
	var aa taws.AwsAccount
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	defer utilsUsualTxFinalize(&tx, &err, &logger, "check-cost")

	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else if aa, err = taws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else if user, err := models.UserByID(db.Db, aa.UserId); err != nil || user.AccountType != "trackit" {
		if err == nil {
			logger.Info("Task 'CheckCost' has been skipped because the user has the wrong account type.", map[string]interface{}{
				"userAccountType": user.AccountType,
				"requiredAccount": "trackit",
			})
		}
	} else {
		runCostCheckForAccount(ctx, aa)
	}
	if err != nil {
		logger.Error("Failed to check account cost.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
	}
	return
}

func getCostFromExplorer(ctx context.Context, aa taws.AwsAccount, intervalBegin, currentMonthBeginning time.Time) (*costexplorer.GetCostAndUsageOutput, error) {
	creds, err := taws.GetTemporaryCredentials(aa, "trackit-cost-check")
	if err != nil {
		return nil, err
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
	}))
	svc := costexplorer.New(sess)

	explorerGranularity := "MONTHLY"
	explorerIntervalStart := intervalBegin.Format("2006-01-02")
	explorerIntervalEnd := currentMonthBeginning.Format("2006-01-02")
	explorerMetrics := "UnblendedCost"
	costExplorerInput := costexplorer.GetCostAndUsageInput{
		Granularity: &explorerGranularity,
		TimePeriod: &costexplorer.DateInterval{
			Start: &explorerIntervalStart,
			End:   &explorerIntervalEnd,
		},
		Metrics: []*string{&explorerMetrics},
	}
	return svc.GetCostAndUsage(&costExplorerInput)
}

func getCostFromES(ctx context.Context, aa taws.AwsAccount, intervalBegin, intervalEnd time.Time) (es.SimplifiedCostsDocument, int, error) {
	accountList := []string{aa.AwsIdentity}
	aggregationParams := []string{"month"}
	indexList := []string{es.IndexNameForUserId(aa.UserId, s3.IndexPrefixLineItem)}
	params := costs.EsQueryParams{
		DateBegin:         intervalBegin,
		DateEnd:           intervalEnd,
		AccountList:       accountList,
		IndexList:         indexList,
		AggregationParams: aggregationParams,
	}
	return costs.MakeElasticSearchRequestAndParseIt(ctx, params)
}

// runCostCheckForAccount runs a cost check for an account
func runCostCheckForAccount(ctx context.Context, aa taws.AwsAccount) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	now := time.Now().UTC()
	currentMonthBeginning := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	intervalEnd := currentMonthBeginning.Add(time.Nanosecond * -1)
	intervalBegin := currentMonthBeginning.AddDate(0, -6, 0)

	explorerCosts, err := getCostFromExplorer(ctx, aa, intervalBegin, currentMonthBeginning)
	if err != nil {
		logger.Error("Failed to retrieve cost from the cost explorer", err.Error())
		return
	}
	esCosts, _, err := getCostFromES(ctx, aa, intervalBegin, intervalEnd)
	if err != nil {
		logger.Error("Failed to retrieve cost from ES", err.Error())
		return
	}
	if len(explorerCosts.ResultsByTime) != len(esCosts.Children) {
		logger.Error("Number of month does not match for cost comparison", map[string]interface{}{
			"explorerMonthNb": len(explorerCosts.ResultsByTime),
			"esMonthNb":       len(esCosts.Children),
		})
		return
	}

	esMonthlyCosts := map[string]int{}

	for _, esMonth := range esCosts.Children {
		esMonthlyCosts[strings.Split(esMonth.Key, "T")[0]] = int(esMonth.Value)
	}

	for _, explorerMonth := range explorerCosts.ResultsByTime {
		explorerMonthCostPtr := explorerMonth.Total["UnblendedCost"].Amount
		if explorerMonthCostPtr == nil {
			logger.Error("Failed to get unblended cost from explorer", nil)
			return
		}
		explorerMonthCost, err := strconv.ParseFloat(*explorerMonthCostPtr, 64)
		if err != nil {
			logger.Error("Failed to parse cost from explorer", err.Error())
			return
		}
		if val, ok := esMonthlyCosts[*explorerMonth.TimePeriod.Start]; ok {
			if val != int(explorerMonthCost) {
				logger.Error("ES cost does not match explorer cost", map[string]interface{}{
					"month":        *explorerMonth.TimePeriod.Start,
					"esCost":       val,
					"explorerCost": int(explorerMonthCost),
				})
				return
			}
		} else {
			logger.Error("Month not found in ES result", map[string]interface{}{
				"month": *explorerMonth.TimePeriod.Start,
			})
			return
		}
	}
	logger.Info("ES costs matches explorer costs over the last 6 months", nil)
}
