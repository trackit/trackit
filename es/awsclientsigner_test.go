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
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

var (
	rightEndpoint   = "https://search-job-msol-prod-trackit-j6ofkgxgmxezkamcywpmqwfsn4.us-west-2.es.amazonaws.com"
	wrongEndpoint   = "https://search-job-msol-prod-trackit-j6ofkgxgmxezkamcywpmqwfsn4.us-west-2.es.amazonaws-com"
	unknownEndpoint = "https://search-job-msol-prod-markets-4nsfwqmpwtcmakzexmgxgkfo6j.us-east-1.es.amazonaws.com"
)

func TestNonExistingCredentials(t *testing.T) {
	_, err := NewSignedElasticClient(rightEndpoint, credentials.AnonymousCredentials)
	if err == nil {
		t.Error("Error should not be nil, instead is nil.")
	}
}

func TestUnknownEndpoint(t *testing.T) {
	_, err := NewSignedElasticClient(unknownEndpoint, credentials.NewSharedCredentials("", ""))
	if err == nil {
		t.Error("Error should not be nil, instead is nil.")
	}
}

func TestWrongEndpoint(t *testing.T) {
	_, err := NewSignedElasticClient(wrongEndpoint, credentials.NewSharedCredentials("", ""))
	if err == nil {
		t.Error("Error should not be nil, instead of nil.")
	} else if err != errWrongEndPoint {
		t.Errorf("Error should be \"%s\", instead is \"%s\".", err.Error(), errWrongEndPoint.Error())
	}
}

func TestRightParameters(t *testing.T) {
	client, err := NewSignedElasticClient(rightEndpoint, credentials.NewSharedCredentials("", ""))
	if err != nil {
		t.Errorf("Error should be nil, instead of \"%s\".", err.Error())
	}
	_, err = client.IndexExists("inexistantindex").Do(context.Background())
	if err != nil {
		t.Errorf("Error should be nil, instead of \"%s\".", err.Error())
	}
}
