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

package route53

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/aws/usageReports"
)

// getRoute53Tags formats []*HostedZone.Tag to map[string]string
func getRoute53Tags(ctx context.Context, hostedZone *route53.HostedZone, svc *route53.Route53) []utils.Tag {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	res := make([]utils.Tag, 0)
	var tagResourceType = route53.TagResourceTypeHostedzone
	tags, err := svc.ListTagsForResource(&route53.ListTagsForResourceInput{
		ResourceId: hostedZone.Id,
		ResourceType: &tagResourceType,
	})

	if err != nil {
		logger.Error("Failed to get Route53 tags", err.Error())
		return res
	}

	for _, tag := range tags.ResourceTagSet.Tags {
		res = append(res, utils.Tag{
			Key:   aws.StringValue(tag.Key),
			Value: aws.StringValue(tag.Value),
		})
	}
	return res
}