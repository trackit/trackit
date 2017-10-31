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

package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/trackit/jsonlog"
)

// AwsTemporaryCredentials hold temporary credentials and various useful data
// in order to assume a role on a client's AWS account.
type AwsTemporaryCredentials struct {
	Expires     time.Time       `json:"expires"`
	Renews      time.Time       `json:"renews"`
	Credentials sts.Credentials `json:"credentials"`
	Account     AwsAccount      `json:"account"`
	SessionName string          `json:"sessionName"`
}

// GetTemporaryCredentials gets temporary credentials in a client's AWS account
// using the STS AssumeRole feature. The returned credentials will last no more
// than an hour. The returned credentials are valid iff the error is nil.
func GetTemporaryCredentials(ctx context.Context, account AwsAccount, sessionName string) (AwsTemporaryCredentials, error) {
	durationSeconds := int64(assumeRoleDuration)
	input := sts.AssumeRoleInput{
		DurationSeconds: &durationSeconds,
		ExternalId:      &account.External,
		RoleArn:         &account.RoleArn,
		RoleSessionName: &sessionName,
	}
	temporaryCredentials := AwsTemporaryCredentials{
		Account:     account,
		SessionName: sessionName,
	}
	var err error
	if result, err := stsService.AssumeRoleWithContext(ctx, &input); err == nil {
		populateTemporaryCredentials(&temporaryCredentials, result)
	} else {
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		logger.Error("Failed to get temporary credentials", map[string]interface{}{
			"error":   err.Error(),
			"account": account,
		})
	}
	return temporaryCredentials, err
}

// populateTemporaryCredentials populates an instance of
// AwsTemporaryCredentials using the results of an sts.AssumeRoleOutput.
func populateTemporaryCredentials(temporaryCredentials *AwsTemporaryCredentials, apiResult *sts.AssumeRoleOutput) {
	temporaryCredentials.Expires = *apiResult.Credentials.Expiration
	temporaryCredentials.Renews = temporaryCredentials.Expires.Add(assumeRoleDuration / 2 * time.Second)
	temporaryCredentials.Credentials = *apiResult.Credentials
}
