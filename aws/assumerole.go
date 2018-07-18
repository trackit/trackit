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
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"

	"github.com/trackit/trackit-server/awsSession"
)

// GetTemporaryCredentials gets temporary credentials in a client's AWS account
// using the STS AssumeRole feature. The returned credentials will last no more
// than an hour. The returned credentials are valid iff the error is nil.
func GetTemporaryCredentials(aa AwsAccount, sessionName string) (*credentials.Credentials, error) {
	creds := stscreds.NewCredentials(awsSession.Session, aa.RoleArn, func(arp *stscreds.AssumeRoleProvider) {
		arp.ExternalID = &aa.External
	})
	_, err := creds.Get()
	if err != nil {
		return nil, err
	} else {
		return creds, nil
	}
}
