package utils

import (
	"fmt"
)

const urlFormat = "https://console.aws.amazon.com/cloudformation/home?region=%s#/stacks/stackinfo?stackId=%s"
const stackId = "aws:cloudformation:stack-id"

func GenCloudFormationUrl(resource *TaggingReportDocument) bool {
	for _, tag := range resource.Tags {
		if tag.Key == stackId {
			resource.CloudFormationURL = fmt.Sprintf(urlFormat, resource.Region, tag.Value)
			return true
		}
	}
	return false
}
