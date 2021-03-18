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

package cloudformation

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws/usageReports"
)

// getCloudFormationTags formats []*StackSummary.Tag to map[string]string
func getCloudFormationTags(ctx context.Context, stack *cloudformation.StackSummary, svc *cloudformation.CloudFormation) []utils.Tag {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	res := make([]utils.Tag, 0)
	tags, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: stack.StackName,
	})
	if err != nil {
		logger.Error("Failed to get CloudFormation tags", err.Error())
		return res
	}

	for _, tag := range tags.Stacks[0].Tags {
		res = append(res, utils.Tag{
			Key:   aws.StringValue(tag.Key),
			Value: aws.StringValue(tag.Value),
		})
	}
	return res
}
