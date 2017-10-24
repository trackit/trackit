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
