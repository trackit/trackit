package entitlement

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/marketplaceentitlementservice"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/models"
)

func CheckUserEntitlements(ctx context.Context, db *sql.Tx, userId int) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	customer, err := models.UserByID(db, userId)
	if err != nil {
		logger.Error("Error while getting customer infos", err)
		return err
	} else if customer.AwsCustomerIdentifier == "" {
		logger.Info("No AWS customer identifier", err)
		return updateCustomerEntitlement(db, ctx, userId, false)
	} else {
		err = checkUserEntitlement(ctx, db, customer.AwsCustomerIdentifier, userId)
		if err != nil {
			logger.Error("Error occured while checking user entitlement", err)
			return err
		}
	}
	return nil
}

// getAwsEntitlementConfig returns an AWS config with the region required for entitlement API calls
func getAwsEntitlementConfig() client.ConfigProvider {
	return session.Must(session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		Region:                        aws.String("us-east-1"),
	}))
}

// getUserEntitlement calls getEntitlements function to retrieve specific user entitlement from AWS marketplace.
func getUserEntitlement(ctx context.Context, customerIdentifier string) ([]*marketplaceentitlementservice.Entitlement, error) {
	svc := marketplaceentitlementservice.New(getAwsEntitlementConfig())
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
		logger.Error("Error when checking the AWS token", map[string]interface{}{
			"message": aerr.Message(),
			"err":     aerr.Error(),
		})
		return nil, err
	}
	return result.Entitlements, nil
}

// checkUserEntitlement enables entitlement to be checked.
func checkUserEntitlement(ctx context.Context, db *sql.Tx, cuId string, userId int) error {
	var expirationDate time.Time
	res, err := getUserEntitlement(ctx, cuId)
	if err != nil {
		return err
	}
	for _, key := range res {
		expirationDate = aws.TimeValue(key.ExpirationDate)
	}
	err = checkExpirationDate(expirationDate, ctx, db, userId)
	return err
}

// checkExpirationDate compares expiration date given by AWS to current time.
// According to result, an update is pushed to db.
func checkExpirationDate(expirationDate time.Time, ctx context.Context, db *sql.Tx, userId int) error {
	var err error
	currentTime := time.Now()
	if expirationDate.After(currentTime) {
		err = updateCustomerEntitlement(db, ctx, userId, true)
	} else {
		err = updateCustomerEntitlement(db, ctx, userId, false)
	}
	return err
}

// updateCustomerEntitlement updates aws customer entitlement according to entitlement value.
func updateCustomerEntitlement(db *sql.Tx, ctx context.Context, userId int, entitlementValue bool) error {
	dbUser, err := models.UserByID(db, userId)
	if err != nil {
		return err
	}
	dbUser.AwsCustomerEntitlement = entitlementValue
	err = dbUser.Update(db)
	return err
}
