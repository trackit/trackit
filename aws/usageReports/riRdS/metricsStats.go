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

package riRdS

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/es/indexes/common"
	"github.com/trackit/trackit/es/indexes/rdsRiReports"
)

//getInstanceTags returns an array of tags associated to the RDS reserved instance given as parameter
func getInstanceTags(ctx context.Context, instance *rds.ReservedDBInstance, svc *rds.RDS) []common.Tag {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	desc := rds.ListTagsForResourceInput{
		ResourceName: instance.ReservedDBInstanceArn,
	}
	res, err := svc.ListTagsForResource(&desc)
	if err != nil {
		logger.Error("Failed to get RDS tags", err.Error())
		return []common.Tag{}
	}
	tags := make([]common.Tag, len(res.TagList))
	for i, tag := range res.TagList {
		tags[i] = common.Tag{
			Key:   aws.StringValue(tag.Key),
			Value: aws.StringValue(tag.Value),
		}
	}
	return tags
}

func getRecurringCharges(reservation *rds.ReservedDBInstance) []rdsRiReports.RecurringCharges {
	charges := make([]rdsRiReports.RecurringCharges, len(reservation.RecurringCharges))
	for i, key := range reservation.RecurringCharges {
		charges[i] = rdsRiReports.RecurringCharges{
			Amount:    aws.Float64Value(key.RecurringChargeAmount),
			Frequency: aws.StringValue(key.RecurringChargeFrequency),
		}
	}
	return charges
}
