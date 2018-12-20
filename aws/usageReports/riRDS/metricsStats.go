//   Copyright 2018 MSolution.IO
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

package riRDS

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit-server/aws/usageReports"
)

//getInstanceTags returns an array of tags associated to the RDS instance given as parameter
func getInstanceTags(ctx context.Context, instance *rds.ReservedDBInstance, svc *rds.RDS) []utils.Tag {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	desc := rds.ListTagsForResourceInput{
		ResourceName: instance.ReservedDBInstanceArn,
	}
	res, err := svc.ListTagsForResource(&desc)
	if err != nil {
		logger.Error("Failed to get RDS tags", err.Error())
		return []utils.Tag{}
	}
	tags := make([]utils.Tag, len(res.TagList))
	for i, tag := range res.TagList {
		tags[i] = utils.Tag{
			Key:   aws.StringValue(tag.Key),
			Value: aws.StringValue(tag.Value),
		}
	}
	return tags
}

func getRecurringCharges(reservation *rds.ReservedDBInstance) []RecurringCharges {
	charges := make([]RecurringCharges, len(reservation.RecurringCharges))
	for i, key := range reservation.RecurringCharges {
		charges[i] = RecurringCharges{
			Amount:    aws.Float64Value(key.RecurringChargeAmount),
			Frequency: aws.StringValue(key.RecurringChargeFrequency),
		}
	}
	return charges
}

// getInstanceCPUStats gets the CPU average and the CPU peak from CloudWatch
//func getInstanceCPUStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, error) {
//	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
//		Namespace:  aws.String("AWS/RDS"),
//		MetricName: aws.String("CPUUtilization"),
//		StartTime:  aws.Time(start),
//		EndTime:    aws.Time(end),
//		Period:     aws.Int64(int64(60*60*24) * 31),
//		Statistics: []*string{aws.String("Average"), aws.String("Maximum")},
//		Dimensions: dimensions,
//	})
//	if err != nil {
//		return -1, -1, err
//	} else if len(stats.Datapoints) > 0 {
//		return aws.Float64Value(stats.Datapoints[0].Average), aws.Float64Value(stats.Datapoints[0].Maximum), nil
//	} else {
//		return -1, -1, nil
//	}
//}

// getInstanceFreeSpaceStats gets the free space stats from CloudWatch
//func getInstanceFreeSpaceStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, float64, error) {
//	freeSpace, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
//		Namespace:  aws.String("AWS/RDS"),
//		MetricName: aws.String("FreeStorageSpace"),
//		StartTime:  aws.Time(start),
//		EndTime:    aws.Time(end),
//		Period:     aws.Int64(int64(60*60*24) * 31),
//		Statistics: []*string{aws.String("Minimum"), aws.String("Maximum"), aws.String("Average")},
//		Dimensions: dimensions,
//	})
//	if err != nil {
//		return -1, -1, -1, err
//	} else if len(freeSpace.Datapoints) > 0 {
//		return aws.Float64Value(freeSpace.Datapoints[0].Minimum),
//			aws.Float64Value(freeSpace.Datapoints[0].Maximum),
//			aws.Float64Value(freeSpace.Datapoints[0].Average), nil
//	} else {
//		return -1, -1, -1, nil
//	}
//}

// getInstanceStats gets the instance stats from CloudWatch
//func getInstanceStats(ctx context.Context, instance *rds.DBInstance, sess *session.Session, start, end time.Time) Stats {
//	logger := jsonlog.LoggerFromContextOrDefault(ctx)
//	svc := cloudwatch.New(sess)
//	dimensions := []*cloudwatch.Dimension{{
//		Name:  aws.String("DBInstanceIdentifier"),
//		Value: aws.String(aws.StringValue(instance.DBInstanceIdentifier)),
//	}}
//	var stats Stats
//	var err error
//	stats.Cpu.Average, stats.Cpu.Peak, err = getInstanceCPUStats(svc, dimensions, start, end)
//	if err != nil {
//		logger.Error("Error when fetching CPU stats from CloudWatch", err.Error())
//	}
//	stats.FreeSpace.Minimum, stats.FreeSpace.Maximum,
//		stats.FreeSpace.Average, err = getInstanceFreeSpaceStats(svc, dimensions, start, end)
//	if err != nil {
//		logger.Error("Error when fetching IO stats from CloudWatch", err.Error())
//	}
//	return stats
//}
