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
	"flag"
	"strconv"
	"time"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/marketplaceentitlementservice"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/awsSession"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
)

// taskCheckEntitlement checks the user Entitlement for AWS Marketplace users
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
		customer, err := models.UserByID(db.Db, userId)
		if err != nil {
			logger.Error("Error while getting cursomer infos", err)
			return err
		} else if customer.AwsCustomerIdentifier == "" {
			logger.Info("No AWS customer identifier", err)
			return nil
		} else {
			err = checkUserEntitlement(ctx, customer.AwsCustomerIdentifier, userId)
			if err != nil {
				logger.Error("Error occured while checking user entitlement", err)
				return err
			}
		}
	}
	return nil
}

// getUserEntitlement calls getEntitlements function to retrieve specific user entitlement from AWS marketplace.
func getUserEntitlement(ctx context.Context, customerIdentifier string, productCode string) ([]*marketplaceentitlementservice.Entitlement, error) {
	svc := marketplaceentitlementservice.New(awsSession.Session)
	var awsInput marketplaceentitlementservice.GetEntitlementsInput
	var filter = make(map[string][]*string)
	filter["CUSTOMER_IDENTIFIER"] = []*string{aws.String(customerIdentifier)}
	awsInput.SetProductCode(config.MarketPlaceProductCode)
	awsInput.SetFilter(filter)
	result, err := svc.GetEntitlements(&awsInput)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !ok {
			return nil, errors.New("AWS error cast failed")
		}
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Error when checking the AWS token", aerr.Message())
		return nil, err
	}
	return result.Entitlements, nil
}

// checkUserEntitlement enables entitlement to be checked.
func checkUserEntitlement(ctx context.Context, cuId string, userId int) error {
	var expirationDate time.Time
	var entitlements []bool
	var err error
	productCodes := strings.Split(config.MarketPlaceProductCode, ":")
	for _, productCode := range productCodes {
		res, err := getUserEntitlement(ctx, cuId, productCode)
		if err != nil {
			return err
		}
		for _, key := range res {
			expirationDate = aws.TimeValue(key.ExpirationDate)
		}
		entitlements = append(entitlements, checkExpirationDate(expirationDate, ctx, db.Db, userId))
	}
	for _, entitlement := range entitlements {
		if entitlement == true {
			err = updateCustomerEntitlement(db.Db, ctx, userId, true)
			return err
		}
	}
	err = updateCustomerEntitlement(db.Db, ctx, userId, false)
	return err
}

// checkExpirationDate compares expiration date given by AWS to current time.
// Return true or false according to the result.
func checkExpirationDate(expirationDate time.Time, ctx context.Context, db *sql.DB, userId int) bool {
	currentTime := time.Now()
	if expirationDate.After(currentTime) {
		return true
	} else {
		return false
	}
}

// updateCustomerEntitlement updates aws customer entitlement according to entitlement value.
func updateCustomerEntitlement(db *sql.DB, ctx context.Context, userId int, entitlementValue bool) error {
	dbUser, err := models.UserByID(db, userId)
	if err != nil {
		return err
	}
	dbUser.AwsCustomerEntitlement = entitlementValue
	err = dbUser.Update(db)
	return err
}
