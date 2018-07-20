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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit-server/awsSession"
	"github.com/trackit/trackit-server/config"
)

var Client *elastic.Client

const (
	retryCount   = 15
	retrySeconds = 2
)

func init() {
	var err error
	logger := jsonlog.DefaultLogger
	options := getElasticSearchConfig()
	for r := retryCount; r > 0; r-- {
		Client, err = elastic.NewClient(options...)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to connect to ElasticSearch database. Retrying in %s seconds.", retrySeconds), err.Error())
			time.Sleep(retrySeconds * time.Second)
		} else {
			logger.Info("Successfully connected to ElasticSearch database.", nil)
			return
		}
	}
	logger.Error("Failed to connect to ElasticSearch database. Not retrying.", nil)
}

// getElasticSearchConfig retrieves the elastic.ClientOptionFunc required to
// correctly configure the server's ElasticSearch client.
func getElasticSearchConfig() []elastic.ClientOptionFunc {
	return []elastic.ClientOptionFunc{
		getElasticSearchUrlConfig(),
		getElasticSearchAuthConfig(),
	}
}

// getElasticSearchUrlConfig gets the SetURL elastic.ClientOptionFunc.
func getElasticSearchUrlConfig() elastic.ClientOptionFunc {
	return elastic.SetURL(config.EsAddress)
}

// getElasticSearchAuthConfig gets the elastic.ClientOptionFunc responsible for
// the authentication of requests on the ElasticSearch server, as necessary.
func getElasticSearchAuthConfig() elastic.ClientOptionFunc {
	authType, authValue := getElasticSearchAuthTypeAndValue()
	logger := jsonlog.DefaultLogger
	switch authType {
	case "basic":
		logger.Debug("Configuring ElasticSearch client with basic auth.", nil)
		return getElasticSearchBasicAuth(authValue)
	case "iam":
		logger.Debug("Configuring ElasticSearch client with IAM auth.", nil)
		return getElasticSearchIamAuth(authValue)
	case "none":
		logger.Debug("Configuring ElasticSearch client with null auth.", nil)
	default:
		logger.Error("Could not configure ElasticSearch client auth: bad auth format.", nil)
		os.Exit(1)
	}
	return configNoop
}

// getElasticSearch gets the configuration required to use basic HTTP
// authentication on the ElasticSearch server.
func getElasticSearchBasicAuth(auth string) elastic.ClientOptionFunc {
	parts := strings.SplitN(auth, ":", 2)
	if len(parts) == 2 {
		return elastic.SetBasicAuth(parts[0], parts[1])
	} else {
		logger := jsonlog.DefaultLogger
		logger.Error("Could not configure ElasticSearch client basic auth: missing username or password.", nil)
		os.Exit(1)
		return nil
	}
}

// getElasticSearchAuthTypeAndValue separates the authentication type from its
// configuration values. The type and values are separated by a colon.
func getElasticSearchAuthTypeAndValue() (string, string) {
	parts := strings.SplitN(config.EsAuthentication, ":", 2)
	if len(parts) == 0 {
		return "none", ""
	} else if len(parts) == 1 {
		return parts[0], ""
	} else {
		return parts[0], parts[1]
	}
}

// getElasticSearchIamAuth gets the type of IAM authentication. Currently only
// EC2Role-based authentication is supported through the "ec2role" value.
func getElasticSearchIamAuth(auth string) elastic.ClientOptionFunc {
	logger := jsonlog.DefaultLogger
	if auth == "ec2role" {
		return getElasticSearchEc2RoleAuth()
	} else {
		logger.Error("Could not configure ElasticSearch client IAM auth: bad value.", nil)
		os.Exit(1)
		return nil
	}
}

// getElasticSearchEc2RoleAuth gets the options to perform AWS v4 signature
// requests to the ElasticSearch server.
func getElasticSearchEc2RoleAuth() elastic.ClientOptionFunc {
	var err error
	if creds := ec2rolecreds.NewCredentials(awsSession.Session); creds != nil {
		if _, err = creds.Get(); err == nil {
			return getElasticSearchEc2RoleAuthOptionFunc(creds)
		}
	} else {
		err = errors.New("got nil credentials")
	}
	jsonlog.DefaultLogger.Error(
		"Could not configure ElasticSearch client IAM auth: failed to retrieve credentials.",
		err.Error(),
	)
	os.Exit(1)
	return nil
}

// getElasticSearchEc2RoleAuthOptionFunc builds the option funcs to sign
// requests with the provided AWS credentials.
func getElasticSearchEc2RoleAuthOptionFunc(creds *credentials.Credentials) elastic.ClientOptionFunc {
	if cofs, err := NewSignedElasticClientOptions(config.EsAddress, creds); err == nil {
		return configEach(cofs...)
	} else {
		jsonlog.DefaultLogger.Error(
			"Could not configure ElasticSearch client IAM auth: failed to create signing HTTP client.",
			err.Error(),
		)
		os.Exit(1)
		return nil
	}
}

// configNoop does not alter the ElasticSearch configuration.
func configNoop(_ *elastic.Client) error { return nil }

// configEach builds an option func which applies all provided option funcs
// sequentially.
func configEach(cofs ...elastic.ClientOptionFunc) elastic.ClientOptionFunc {
	return func(c *elastic.Client) error {
		for _, cof := range cofs {
			if err := cof(c); err != nil {
				return err
			}
		}
		return nil
	}
}
