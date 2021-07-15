package utils

import (
	"fmt"
)

const urlFormat = "https://console.aws.amazon.com/cloudformation/home?region=%s#/stacks/stackinfo?stackId=%s"
const stackId = "aws:cloudformation:stack-id"

func (doc TaggingReportDocument) GenCloudFormationUrl() TaggingReportDocument {
	for _, tag := range doc.Tags {
		if tag.Key == stackId && tag.Value != "" {
			doc.CloudFormationURL = fmt.Sprintf(urlFormat, doc.Region, tag.Value)
			return doc
		}
	}
	return doc
}
