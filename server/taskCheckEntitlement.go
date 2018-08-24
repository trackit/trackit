//   Copyright 2018 MSolution.IO
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
	"flag"
	"strconv"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/db"
	//"github.com/aws/aws-sdk-go/service/marketplacemetering"
	"github.com/aws/aws-sdk-go/service/marketplaceentitlementservice"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/trackit/trackit-server/config"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

// taskCheckEntitlement check the user Entitlement for AWS Marketplace users
func taskCheckEntitlement(ctx context.Context) error {
	args := flag.Args()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'checkuserentitlement'.", map[string]interface{}{
		"args": args,
	})
	if len(args) != 1 {
		return errors.New("taskCheckEntitlement requires one integer argument")
	} else if userId, err := strconv.Atoi(args[0]); err != nil {
		return err
	} else {
		cuId, error := getCustomerIdentifier(db.Db, ctx, userId)
		if error != nil {
			return error
		} else if cuId == "" {
			return nil
		} else {
			getUserEntitlement(ctx, cuId)
			//res := getUserEntitlement(ctx, cuId)
			//if res["Entitlement"] == "Ok" {
			//	updateCustomerEntitlement(db.Db, ctx, userId, 0)
			//} else {
			//	updateCustomerEntitlement(db.Db, ctx, userId, 1)
			//}
		}
	}
	return nil
}

func getUserEntitlement(ctx context.Context, customerIdentifier string) (*GetEntitlementsOutput, error){
	mySession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion),
	}))
	svc := marketplaceentitlementservice.New(mySession)
	var awsInput marketplaceentitlementservice.GetEntitlementsInput
	var filter = make(map[string][]*string)
	filter["CUSTOMER_IDENTIFIER"] = make([]*string, 0)
	filter["CUSTOMER_IDENTIFIER"] = append(filter["CUSTOMER_IDENTIFIER"], &customerIdentifier)
	awsInput.SetProductCode("f2k4zzzgvihcmlgrf4qlyo4sr")
	awsInput.SetFilter(filter)
	result, err := svc.GetEntitlements(&awsInput)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !ok {
			return nil, errors.New("AWS error cast failed")
		}
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Error when checking the AWS token", aerr.Message())
		return nil, nil
	}
	fmt.Print("AWS Customer Entitlement Result : ")
	fmt.Print(result.Entitlements)
	return result, nil
}

func getCustomerIdentifier(db *sql.DB, ctx context.Context, userId int) (string, error) {
	const sqlstr = `SELECT aws_customer_identifier FROM user WHERE id = ?`
	res, err := db.Query(sqlstr, userId)
	defer res.Close()
	if err != nil {
		return "", err
	}
	var token string
	for res.Next() {
		err := res.Scan(&token)
		if err != nil {
			return "", err
		}
	}
	return token, nil
}

func updateCustomerEntitlement(db *sql.DB, ctx context.Context, userId int, entitlementValue int) (error) {
	const sqlstr = `UPDATE user SET
		aws_customer_entitlement=?,
	WHERE id=?`
		var err error = nil
		//if err != nil {
		//	errorValue = err.Error()
		//}
		_, err = db.Exec(sqlstr, entitlementValue, userId)
		return err
}

//func checkUserEntitlement(ctx context.Context, amId int, cuId string) {
//	var entitlementInput marketplaceentitlementservice.GetEntitlementsInput
//	mySession := session.Must(session.NewSession(&aws.Config{
//		Region: aws.String(config.AwsRegion),
//	}))
//	svc := marketplaceentitlementservice.New(mySession)
//	entitlementInput.SetProductCode("f2k4zzzgvihcmlgrf4qlyo4sr")
//	entitlementInput.SetFilter(map["CUSTOMER_IDENTIFIER"][cuId])
//}
//
//func EntitlementUpdate(db *sql.DB, br s3.BillRepository) (int64, error) {
//	const sqlstr = `INSERT INTO user(
//		aws_entitlement
//	) VALUES (?)`
//	res, err := db.Exec(sqlstr, br.Id, backendId)
//	if err != nil {
//		return 0, err
//	}
//	return res.LastInsertId()
//}
//
//func updateCompletion(ctx context.Context, aaId, brId int, db *sql.DB, updateId int64, err error) {
//	rErr := registerUpdateCompletion(db, updateId, err)
//	if rErr != nil {
//		logger := jsonlog.LoggerFromContextOrDefault(ctx)
//		logger.Error("Failed to register ingestion completion.", map[string]interface{}{
//			"awsAccountId":     aaId,
//			"billRepositoryId": brId,
//			"error":            rErr.Error(),
//			"updateId":         updateId,
//		})
//	}
//}
//
//func registerUpdateCompletion(db *sql.DB, updateId int64, err error) error {
//	const sqlstr = `UPDATE aws_bill_update_job SET
//		completed=?,
//		error=?
//	WHERE id=?`
//	var errorValue string
//	var now = time.Now()
//	if err != nil {
//		errorValue = err.Error()
//	}
//	_, err = db.Exec(sqlstr, now, errorValue, updateId)
//	return err
//}
//
//const (
//	UpdateIntervalMinutes = 6 * 60
//	UpdateIntervalWindow  = 2 * 60
//)
//
//// updateBillRepositoryForNextUpdate plans the next update for a
//// BillRepository.
//func updateBillRepositoryForNextUpdate(ctx context.Context, tx *sql.Tx, br s3.BillRepository, latestManifest time.Time) error {
//	if latestManifest.After(br.LastImportedManifest) {
//		br.LastImportedManifest = latestManifest
//	}
//	updateDeltaMinutes := time.Duration(UpdateIntervalMinutes-UpdateIntervalWindow/2+rand.Int63n(UpdateIntervalWindow)) * time.Minute
//	br.NextUpdate = time.Now().Add(updateDeltaMinutes)
//	return s3.UpdateBillRepository(br, tx)
//}
