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
	"github.com/trackit/trackit/es/indexes/common"
)

// fetchOneRegion establishes a new API connection for the corresponding region and gets all the tags for all the buckets associated with it (passed through region_buckets
func fetchOneRegion(ctx context.Context, creds *credentials.Credentials, region string, region_buckets []s3.Bucket, bucketChan chan Bucket) {
	defer close(bucketChan)
	session := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	service := s3.New(session)
	for _, bucket := range region_buckets {
		bucketChan <- Bucket{
			BucketBase: BucketBase{
				Name:         aws.StringValue(bucket.Name),
				CreationDate: aws.TimeValue(bucket.CreationDate),
				Region:       region,
			},
			Tags: getS3Tags(ctx, &bucket, service),
		}
	}
}

// getRegionBucketMap creates a map of regions (as strings) and all the buckets associated with it
func getRegionBucketMap(ctx context.Context, creds *credentials.Credentials, sess *session.Session) (map[string][]s3.Bucket, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := s3.New(sess)
	buckets, err := svc.ListBuckets(nil)
	if err != nil {
		logger.Error("Error when describing S3 Buckets", err.Error())
		return nil, err
	}

	region_map := make(map[string][]s3.Bucket)
	for _, bucket := range buckets.Buckets {
		locationRes, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			logger.Warning("Failed to get bucket location", err.Error())
			continue
		}
		key := aws.StringValue(locationRes.LocationConstraint)
		if key == "" {
			key = "us-east-1"
		}
		region_map[key] = append(region_map[key], *bucket)
	}
	return region_map, nil
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

	// You might find yourself thinking this can be eluded because S3 buckets are supposed to be global, meaning it should be possible to just open one connection to one region and do everything from there. Do not delude yourself, for this is the path to learning that GetBucketTagging requires a connection to the region corresponding to the bucket, which is why we need to do this
	regions, err := getRegionBucketMap(ctx, creds, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return err
	}
	bucketChans := make([]<-chan Bucket, 0, len(regions))
	for regionName, regionBuckets := range regions {
		bucketChan := make(chan Bucket)
		go fetchOneRegion(ctx, creds, regionName, regionBuckets, bucketChan)
		bucketChans = append(bucketChans, bucketChan)
	}
	buckets := make([]BucketReport, 0)
	for bucket := range merge(bucketChans...) {
		buckets = append(buckets, BucketReport{
			ReportBase: common.ReportBase{
				Account:    account,
				ReportDate: now,
				ReportType: "daily",
			},
			Bucket: bucket,
		})
	}
	return importS3ToEs(ctx, awsAccount, buckets)
}
