package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
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
func GetTemporaryCredentials(account AwsAccount, sessionName string) (AwsTemporaryCredentials, error) {
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
	if result, err := stsService.AssumeRole(&input); err == nil {
		populateTemporaryCredentials(&temporaryCredentials, result)
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
