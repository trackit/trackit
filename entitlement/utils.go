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

package entitlement

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/marketplaceentitlementservice"
	"github.com/trackit/jsonlog"
)

// getAwsEntitlementConfig returns an AWS config with the region required for entitlement API calls
func getAwsEntitlementConfig() client.ConfigProvider {
	return session.Must(session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		Region:                        aws.String("us-east-1"),
	}))
}

// getAwsEntitlement calls getEntitlements function to retrieve specific user entitlement from AWS marketplace.
func getAwsEntitlement(ctx context.Context, customerIdentifier string, productCode string) (*marketplaceentitlementservice.Entitlement, error) {
	svc := marketplaceentitlementservice.New(getAwsEntitlementConfig())
	var awsInput marketplaceentitlementservice.GetEntitlementsInput
	var filter = make(map[string][]*string)
	filter["CUSTOMER_IDENTIFIER"] = []*string{aws.String(customerIdentifier)}
	awsInput.SetProductCode(productCode)
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
	if len(result.Entitlements) == 0 {
		return nil, nil
	}
	return result.Entitlements[0], nil
}
