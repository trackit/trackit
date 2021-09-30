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

// Package plugins_utils implements common utilities to help with implementing plugins
package plugins_utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetEc2ClientSession is a utility function to create an ec2 session
// it takes credentials and a region and returns an ec2 session
func GetEc2ClientSession(creds *credentials.Credentials, region *string) *ec2.EC2 {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      region,
	}))
	return ec2.New(sess)
}

// GetS3ClientSession is a utility function to create an S3 session
// it takes credentials returns an S3 session
func GetS3ClientSession(creds *credentials.Credentials) *s3.S3 {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
	}))
	return s3.New(sess)
}
