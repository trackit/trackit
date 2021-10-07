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

package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws/usageReports"
)

// getSQSTags formats []*queue.Tag to map[string]string
func getSQSTags(ctx context.Context, queueUrl *string, svc *sqs.SQS) []utils.Tag {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	res := make([]utils.Tag, 0)
	tags, err := svc.ListQueueTags(&sqs.ListQueueTagsInput{
		QueueUrl: queueUrl,
	})
	if err != nil {
		logger.Error("Failed to get SQS tags", err.Error())
		return res
	}

	for keyTag, valueTag := range tags.Tags {
		res = append(res, utils.Tag{
			Key:   keyTag,
			Value: aws.StringValue(valueTag),
		})
	}
	return res
}
