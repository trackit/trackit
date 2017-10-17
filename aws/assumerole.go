package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
)

type AwsTemporaryCredentials struct {
	Expires     time.Time       `json:"expires"`
	Renews      time.Time       `json:"renews"`
	Credentials sts.Credentials `json:"credentials"`
	Account     AwsAccount      `json:"account"`
	SessionName string          `json:"sessionName"`
}

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

func populateTemporaryCredentials(temporaryCredentials *AwsTemporaryCredentials, apiResult *sts.AssumeRoleOutput) {
	temporaryCredentials.Expires = *apiResult.Credentials.Expiration
	temporaryCredentials.Renews = temporaryCredentials.Expires.Add(assumeRoleDuration / 2 * time.Second)
	temporaryCredentials.Credentials = *apiResult.Credentials
}
