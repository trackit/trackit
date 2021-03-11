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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/aws/usageReports"
)

// getS3Tags formats []*bucket.Tag to map[string]string
func getS3Tags(ctx context.Context, bucket *s3.Bucket, svc *s3.S3) []utils.Tag {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	res := make([]utils.Tag, 0)
	tags, err := svc.GetBucketTagging(&s3.GetBucketTaggingInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		logger.Error("Failed to get S3 tags", err.Error())
		return res
	}

	for _, tag := range tags.TagSet {
		res = append(res, utils.Tag{
			Key:   aws.StringValue(tag.Key),
			Value: aws.StringValue(tag.Value),
		})
	}
	return res
}