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
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/olivere/elastic"
	"github.com/sha1sum/aws_signing_client"
)

const (
	// The URI contains 5 mendatory parts split by '.'
	// domainname.region.service.amazonaws.com
	endpointMetaDataLengthRequirement = 5
)

var (
	errWrongEndPoint = errors.New("Wrong endpoint parameter")
)

// checkParametersError checks errors from NewSignedElasticClient's parameters.
func checkParametersError(endpointMetaData []string, creds *credentials.Credentials) error {
	if _, err := creds.Get(); err != nil {
		return err
	} else if len(endpointMetaData) < endpointMetaDataLengthRequirement {
		return errWrongEndPoint
	}
	return nil
}

// NewSignedElasticClient creates a signed *elastic.Client ready for using with AWS ElasticSearch.
// It takes as parameter:
//		- endpoint: The endpoint URI gettable from AWS.
//		- creds: Credentials from AWS/Credentials.
func NewSignedElasticClient(endpoint string, creds *credentials.Credentials) (*elastic.Client, error) {
	if cofs, err := NewSignedElasticClientOptions(endpoint, creds); err == nil {
		cof := configEach(cofs...)
		if ec, err := elastic.NewClient(elastic.SetURL(endpoint), cof); err == nil {
			return ec, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// NewSignedElasticClientOptions builds elastic client option funcs which
// configure an ElasticSearch client to use AWSv4 signature.
func NewSignedElasticClientOptions(endpoint string, creds *credentials.Credentials) ([]elastic.ClientOptionFunc, error) {
	if httpClient, err := NewSignedHttpClientForElasticSearch(endpoint, creds); err != nil {
		return nil, err
	} else {
		return []elastic.ClientOptionFunc{
			elastic.SetScheme("https"),
			elastic.SetHttpClient(httpClient),
			elastic.SetSniff(false),
		}, nil
	}
}

// NewSignedHttpClientForElasticSearch returns an http.Client which signs its
// requests with AWS v4 signatures, for ElasticSearch only.
func NewSignedHttpClientForElasticSearch(endpoint string, creds *credentials.Credentials) (*http.Client, error) {
	endpointParts := strings.Split(endpoint, ".")
	if err := checkParametersError(endpointParts, creds); err != nil {
		return nil, err
	}
	region := endpointParts[len(endpointParts)-4]
	return NewSignedHttpClient(creds, region, "es")
}

// NewSignedHttpClient returns an http.Client which signs its requests with AWS
// v4 signatures for the provided service name and region.
func NewSignedHttpClient(creds *credentials.Credentials, region, service string) (*http.Client, error) {
	signer := v4.NewSigner(creds)
	return aws_signing_client.New(signer, nil, service, region)
}
