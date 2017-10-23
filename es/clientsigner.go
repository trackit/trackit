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

package es

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/sha1sum/aws_signing_client"
	"gopkg.in/olivere/elastic.v5"
)

const (
	// The URI contains 5 mendatory parts split by '.'
	// domainname.region.service.amazonaws.com
	endpointMetaDataLengthRequirement = 5
)

// checkParametersError checks errors from NewSignedElasticClient's parameters.
func checkParametersError(endpointMetaData []string, creds *credentials.Credentials) error {
	if _, err := creds.Get(); err != nil {
		return err
	} else if len(endpointMetaData) < endpointMetaDataLengthRequirement {
		return errors.New("Wrong endpoint parameter")
	}
	return nil
}

// NewSignedElasticClient creates a signed *elastic.Client ready for using with AWS ElasticSearch.
// It takes as parameter:
//		- endpoint: The endpoint URI gettable from AWS.
//		- creds: Credentials from AWS/Credentials.
func NewSignedElasticClient(endpoint string, creds *credentials.Credentials) (*elastic.Client, error) {
	awsSigner := v4.NewSigner(creds)
	endpointMetaData := strings.Split(endpoint, ".")
	if err := checkParametersError(endpointMetaData, creds); err != nil {
		return nil, err
	}
	awsRegion := endpointMetaData[len(endpointMetaData)-4]
	awsClient, err := aws_signing_client.New(awsSigner, nil, "es", awsRegion)
	if err != nil {
		return nil, err
	}
	prefix := ""
	if !strings.HasPrefix(endpoint, "http") {
		prefix = "https://"
	}
	return elastic.NewClient(
		elastic.SetURL(prefix+endpoint),
		elastic.SetScheme("https"),
		elastic.SetHttpClient(awsClient),
		elastic.SetSniff(false),
	)
}
