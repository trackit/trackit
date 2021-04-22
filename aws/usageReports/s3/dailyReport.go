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

package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
)

// fetchDailyS3List sends in bucketChan the S3 Buckets fetched from ListBuckets
func fetchDailyS3List(ctx context.Context, creds *credentials.Credentials, region string, bucketChan chan Bucket) error {
	defer close(bucketChan)
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := s3.New(sess)
	buckets, err := svc.ListBuckets(nil)
	if err != nil {
		logger.Error("Error when describing S3 Buckets", err.Error())
		return err
	}
	for _, bucket := range buckets.Buckets {
		locationRes, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			continue
		}
		if region == aws.StringValue(locationRes.LocationConstraint) {
			bucketChan <- Bucket{
				BucketBase: BucketBase{
					Name:         aws.StringValue(bucket.Name),
					CreationDate: aws.TimeValue(bucket.CreationDate),
					Region:       region,
				},
				Tags: getS3Tags(ctx, bucket, svc),
			}
		}
	}
	return nil
}

// FetchDailyS3Stats fetches the stats of the S3 Bucket of an AwsAccount
// to import them in ElasticSearch. The stats are fetched from the last hour.
// In this way, FetchDailyS3Stats should be called every hour.
func FetchDailyS3Stats(ctx context.Context, awsAccount taws.AwsAccount) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching S3 stats", map[string]interface{}{"awsAccountId": awsAccount.Id})
	creds, err := taws.GetTemporaryCredentials(awsAccount, MonitorS3StsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	now := time.Now().UTC()
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return err
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return err
	}
	bucketChans := make([]<-chan Bucket, 0, len(regions))
	for _, region := range regions {
		bucketChan := make(chan Bucket)
		go fetchDailyS3List(ctx, creds, region, bucketChan)
		bucketChans = append(bucketChans, bucketChan)
	}
	buckets := make([]BucketReport, 0)
	for bucket := range merge(bucketChans...) {
		buckets = append(buckets, BucketReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Bucket: bucket,
		})
	}
	return importS3ToEs(ctx, awsAccount, buckets)
}
