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
	"strconv"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports/cloudformation"
	"github.com/trackit/trackit/aws/usageReports/ebs"
	"github.com/trackit/trackit/aws/usageReports/ec2"
	"github.com/trackit/trackit/aws/usageReports/elasticache"
	"github.com/trackit/trackit/aws/usageReports/es"
	"github.com/trackit/trackit/aws/usageReports/history"
	"github.com/trackit/trackit/aws/usageReports/lambda"
	"github.com/trackit/trackit/aws/usageReports/rds"
	"github.com/trackit/trackit/aws/usageReports/riEc2"
	"github.com/trackit/trackit/aws/usageReports/riRdS"
	"github.com/trackit/trackit/aws/usageReports/route53"
	"github.com/trackit/trackit/aws/usageReports/s3"
	"github.com/trackit/trackit/aws/usageReports/sqs"
	"github.com/trackit/trackit/aws/usageReports/stepfunction"
	"github.com/trackit/trackit/cache"
	"github.com/trackit/trackit/db"
	onDemandToRiEc2 "github.com/trackit/trackit/onDemandToRI/ec2"
)

const invalidAccId = -1

var invalidTime = time.Time{}

type processAccount func(ctx context.Context, aa aws.AwsAccount) error

type accountProcessor struct {
	ErrName string
	Run     processAccount
}

// Every processor of accountProcessors will be run and will fetch those services.
// Don't forget to add the returned error to the database.
// registerAccountProcessingCompletion will push errors to the aws_account_update_job table
var accountProcessors = []accountProcessor{
	{
		ErrName: "cloudformation",
		Run:     processAccountCloudFormation,
	}, {
		ErrName: "route53",
		Run:     processAccountRoute53,
	}, {
		ErrName: "s3",
		Run:     processAccountS3,
	}, {
		ErrName: "sqs",
		Run:     processAccountSQS,
	}, {
		ErrName: "stepfunction",
		Run:     processAccountStepFunction,
	}, {
		ErrName: "ec2",
		Run:     processAccountEC2,
	}, {
		ErrName: "rds",
		Run:     processAccountRDS,
	}, {
		ErrName: "es",
		Run:     processAccountES,
	}, {
		ErrName: "elasticcache",
		Run:     processAccountElastiCache,
	}, {
		ErrName: "lambda",
		Run:     processAccountLambda,
	}, {
		ErrName: "ec2-ri",
		Run:     riEc2.FetchDailyReservationsStats,
	}, {
		ErrName: "rds-ri",
		Run:     riRdS.FetchDailyInstancesStats,
	}, {
		ErrName: "od-ec2-ri",
		Run:     onDemandToRiEc2.RunOnDemandToRiEc2,
	}, {
		ErrName: "ebs",
		Run:     processAccountEbsSnapshot,
	},
}

func checkArgumentsForProcessAccount(args []string) (int, time.Time, error) {
	if len(args) < 1 {
		return invalidAccId, invalidTime, errors.New("process account task requires at least an integer argument as AWS Account UD")
	}
	accId, err := strconv.Atoi(args[0])
	if err != nil {
		return invalidAccId, invalidTime, err
	}
	now := time.Now().UTC()
	if len(args) == 3 {
		if month, err := strconv.Atoi(args[1]); err != nil {
			return invalidAccId, invalidTime, err
		} else if year, err := strconv.Atoi(args[2]); err != nil {
			return invalidAccId, invalidTime, err
		} else {
			formattedMonth := time.Month(month)
			date := time.Date(year, formattedMonth, 1, 0, 0, 0, 0, now.Location()).UTC()
			return accId, date, nil
		}
	}
	return accId, invalidTime, nil
}

// taskProcessAccount processes an AwsAccount to retrieve data from the AWS api.
func taskProcessAccount(ctx context.Context) error {
	args := paramsFromContextOrArgs(ctx)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'process-account'.", map[string]interface{}{
		"args": args,
	})
	aaId, date, err := checkArgumentsForProcessAccount(args)
	if err != nil {
		return err
	} else {
		return ingestDataForAccount(ctx, aaId, date)
	}
}

// ingestDataForAccount ingests the AWS api data for an AwsAccount.
func ingestDataForAccount(ctx context.Context, aaId int, date time.Time) (err error) {
	var tx *sql.Tx
	var aa aws.AwsAccount
	var updateId int64
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	defer utilsUsualTxFinalize(&tx, &err, &logger, "account data injestion")

	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else if aa, err = aws.GetAwsAccountWithId(aaId, tx); err != nil {
	} else if updateId, err = registerAccountProcessing(db.Db, aa); err != nil {
	} else {
		var processAccountErrors = make(map[string]error)
		if date.IsZero() {
			for _, accountProcessor := range accountProcessors {
				processAccountErrors[accountProcessor.ErrName] = accountProcessor.Run(ctx, aa)
			}
		}
		historyCreated, historyErr := processAccountHistory(ctx, aa, date)
		updateAccountProcessingCompletion(ctx, aaId, db.Db, updateId, nil, processAccountErrors, historyErr, historyCreated)
	}
	if err != nil {
		updateAccountProcessingCompletion(ctx, aaId, db.Db, updateId, err, nil, nil, false)
		logger.Error("Failed to process account data.", map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
		})
	}
	var affectedRoutes = []string{
		"/ec2",
		"/ec2/coverage",
		"/ec2/unused",
		"/elasticache",
		"/elasticache/unused",
		"/es",
		"/es/unused",
		"/lambda",
		"/rds",
		"/rds/unused",
		"/ri/ec2",
		"/ri/rds",
	}
	err = cache.RemoveMatchingCache(affectedRoutes, []string{aa.AwsIdentity}, logger)
	return
}

func registerAccountProcessing(db *sql.DB, aa aws.AwsAccount) (int64, error) {
	const sqlstr = `INSERT INTO aws_account_update_job(
		aws_account_id,
		worker_id
	) VALUES (?, ?)`
	res, err := db.Exec(sqlstr, aa.Id, backendId)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func updateAccountProcessingCompletion(ctx context.Context, aaId int, db *sql.DB, updateId int64, jobErr error, processAccountErrors map[string]error, historyErr error, historyCreated bool) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	doErr := func(errorText string, err error) {
		logger.Error(errorText, map[string]interface{}{
			"awsAccountId": aaId,
			"error":        err.Error(),
			"updateId":     updateId,
		})
	}
	if uErr := updateNextUpdateAccount(db, aaId); uErr != nil {
		doErr("Failed to update the account next update date", uErr)
	}
	if rErr := registerAccountProcessingCompletion(db, updateId, jobErr, processAccountErrors, historyErr, historyCreated); rErr != nil {
		doErr("Failed to register account processing completion", rErr)
	}
}

func errToStr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func updateNextUpdateAccount(db *sql.DB, aaId int) error {
	const sqlstr = `UPDATE aws_account SET
		next_update=?
	WHERE id=?`
	_, err := db.Exec(sqlstr, time.Now().AddDate(0, 0, 1), aaId)
	return err
}

func registerAccountProcessingCompletion(db *sql.DB, updateId int64, jobErr error, errors map[string]error, historyErr error, historyCreated bool) error {
	const sqlstr = `UPDATE aws_account_update_job SET
		completed=?,
		jobError=?,
		rdsError=?,
		ec2Error=?,
		esError=?,
		elastiCacheError=?,
		lambdaError=?,
		ebsError=?,
		riEc2Error=?,
		riRdsError=?,
		odToRiEc2Error=?,
		stepFunctionError=?,
		s3Error=?,
		sqsError=?,
		cloudFormationError=?,
		route53Error=?,
		historyError=?,
		monthly_reports_generated=?
	WHERE id=?`
	_, err := db.Exec(
		sqlstr,
		time.Now(),
		errToStr(jobErr),
		errToStr(errors["rds"]),
		errToStr(errors["ec2"]),
		errToStr(errors["es"]),
		errToStr(errors["elasticcache"]),
		errToStr(errors["lambda"]),
		errToStr(errors["ebs"]),
		errToStr(errors["ec2-ri"]),
		errToStr(errors["eds-ri"]),
		errToStr(errors["ec2-od-ri"]),
		errToStr(errors["stepfunction"]),
		errToStr(errors["s3"]),
		errToStr(errors["sqs"]),
		errToStr(errors["cloudformation"]),
		errToStr(errors["route53"]),
		errToStr(historyErr),
		historyCreated,
		updateId)
	return err
}

// processAccountRDS processes all the RDS data for an AwsAccount
func processAccountRDS(ctx context.Context, aa aws.AwsAccount) error {
	err := rds.FetchDailyInstancesStats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest RDS data.", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountS3 processes all the S3 data for an AwsAccount
func processAccountS3(ctx context.Context, aa aws.AwsAccount) error {
	err := s3.FetchDailyS3Stats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest S3 data.", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountCloudFormation processes all the S3 data for an AwsAccount
func processAccountCloudFormation(ctx context.Context, aa aws.AwsAccount) error {
	err := cloudformation.FetchDailyCloudFormationStats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest CloudFormation data.", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountSQS processes all the SQS data for an AwsAccount
func processAccountSQS(ctx context.Context, aa aws.AwsAccount) error {
	err := sqs.FetchDailySQSStats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest SQS data.", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountEC2 processes all the EC2 data for an AwsAccount
func processAccountEC2(ctx context.Context, aa aws.AwsAccount) error {
	err := ec2.FetchDailyInstancesStats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest EC2 data.", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountES processes all the ES data for an AwsAccount
func processAccountES(ctx context.Context, aa aws.AwsAccount) error {
	err := es.FetchDomainsStats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest ES data", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountRoute53 processes all the StepFunctions data for an AwsAccount
func processAccountRoute53(ctx context.Context, aa aws.AwsAccount) error {
	err := route53.FetchDailyRoute53Stats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest Route53 data", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountStepFunction processes all the StepFunctions data for an AwsAccount
func processAccountStepFunction(ctx context.Context, aa aws.AwsAccount) error {
	err := stepfunction.FetchDailyStepFunctionsStats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest StepFunction data", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountElastiCache processes all the ElastiCache data for an AwsAccount
func processAccountElastiCache(ctx context.Context, aa aws.AwsAccount) error {
	err := elasticache.FetchDailyInstancesStats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest ElastiCache data", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountEbsSnapshot process all the EBS Snapshot data for an AwsAccount
func processAccountEbsSnapshot(ctx context.Context, aa aws.AwsAccount) error {
	err := ebs.FetchDailySnapshotsStats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest EBS Snapshots data", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountLambda processes all the Lambda data for an AwsAccount
func processAccountLambda(ctx context.Context, aa aws.AwsAccount) error {
	err := lambda.FetchDailyFunctionsStats(ctx, aa)
	if err != nil {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest Lambda data", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return err
}

// processAccountHistory processes EC2, RDS, ES, ElastiCache, Lambda and EC2 Coverage data with billing data for an AwsAccount
func processAccountHistory(ctx context.Context, aa aws.AwsAccount, date time.Time) (bool, error) {
	status, err := history.FetchHistoryInfos(ctx, aa, date)
	if err != nil && err != history.ErrBillingDataIncomplete {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to ingest History data.", map[string]interface{}{
			"awsAccountId": aa.Id,
			"error":        err.Error(),
		})
	}
	return status, err
}
