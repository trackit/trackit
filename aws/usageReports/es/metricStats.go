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

package es

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws/usageReports"
)

// getDomainCPUStats gets the CPU average and the CPU peak from CloudWatch
func getDomainCPUStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, error) {
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ES"),
		MetricName: aws.String("CPUUtilization"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30),
		Statistics: []*string{aws.String("Average"), aws.String("Maximum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return -1, -1, err
	} else if len(stats.Datapoints) > 0 {
		return aws.Float64Value(stats.Datapoints[0].Average), aws.Float64Value(stats.Datapoints[0].Maximum), nil
	} else {
		return -1, -1, nil
	}
}

// getDomainFreeStorage gets the domain free storage from CloudWatch
func getDomainFreeStorage(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, error) {
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ES"),
		MetricName: aws.String("FreeStorageSpace"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30),
		Statistics: []*string{aws.String("Average")},
		Dimensions: dimensions,
	})
	if err != nil {
		return -1, err
	} else if len(stats.Datapoints) > 0 {
		return aws.Float64Value(stats.Datapoints[0].Average), nil
	} else {
		return -1, nil
	}
}

// getDomainJVMMemoryPressure gets the domain Maximum JVM Memory Pressure from CloudWatch
func getDomainJVMMemoryPressure(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, error) {
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ES"),
		MetricName: aws.String("JVMMemoryPressure"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 30),
		Statistics: []*string{aws.String("Maximum"), aws.String("Average")},
		Dimensions: dimensions,
	})
	if err != nil {
		return -1, -1, err
	} else if len(stats.Datapoints) > 0 {
		return aws.Float64Value(stats.Datapoints[0].Maximum), aws.Float64Value(stats.Datapoints[0].Average), nil
	} else {
		return -1, -1, nil
	}
}

// getDomainStats gets the domains stats from CloudWatch
func getDomainStats(ctx context.Context, domain string, sess *session.Session, start, end time.Time) Stats {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := cloudwatch.New(sess)
	accountId, err := utils.GetAccountId(ctx, sess)
	if err != nil {
		logger.Error("Error while getting the account id", err.Error())
		return Stats{Cpu{-1, -1}, -1, JVMMemoryPressure{-1, -1}}
	}
	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("DomainName"),
			Value: aws.String(domain),
		},
		{
			Name:  aws.String("ClientId"),
			Value: aws.String(accountId),
		},
	}
	var stats Stats
	stats.Cpu.Average, stats.Cpu.Peak, err = getDomainCPUStats(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching CPU stats from CloudWatch", err.Error())
	}
	stats.FreeSpace, err = getDomainFreeStorage(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching storage stats from CloudWatch", err.Error())
	}
	stats.JVMMemoryPressure.Peak, stats.JVMMemoryPressure.Average, err = getDomainJVMMemoryPressure(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching JVM Memory stats from CloudWatch", err.Error())
	}
	return stats
}
