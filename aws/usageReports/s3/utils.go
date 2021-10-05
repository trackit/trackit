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
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/common"
)

const MonitorS3StsSessionName = "monitor-s3"

type (
	// BucketReport is saved in ES to have all the information of a S3 Bucket
	BucketReport struct {
		common.ReportBase
		Bucket Bucket `json:"bucket"`
	}

	// BucketBase contains basics information of a S3 Bucket
	BucketBase struct {
		Name         string    `json:"name"`
		CreationDate time.Time `json:"creationDate"`
		Region       string    `json:"region"`
	}

	// Bucket contains all the information of an S3 Bucket
	Bucket struct {
		BucketBase
		Tags []common.Tag `json:"tags"`
	}
)

// importS3ToEs imports S3 in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importS3ToEs(ctx context.Context, aa taws.AwsAccount, buckets []BucketReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating S3 for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixS3Report)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, bucket := range buckets {
		id, err := generateId(bucket)
		if err != nil {
			logger.Error("Error when marshaling bucket var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, bucket, TypeS3Report, index, id)
	}
	err = bp.Flush()
	if closeErr := bp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		logger.Error("Fail to put S3 in ES", err.Error())
		return err
	}
	logger.Info("S3 put in ES", nil)
	return nil
}

func generateId(bucket BucketReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
		Type       string    `json:"reportType"`
	}{
		bucket.Account,
		bucket.ReportDate,
		bucket.Bucket.Name,
		bucket.ReportType,
	})
	if err != nil {
		return "", err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	return hash64, nil
}

// merge function from https://blog.golang.org/pipelines#TOC_4
// It allows to merge many chans to one.
func merge(cs ...<-chan Bucket) <-chan Bucket {
	var wg sync.WaitGroup
	out := make(chan Bucket)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Bucket) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done. This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
