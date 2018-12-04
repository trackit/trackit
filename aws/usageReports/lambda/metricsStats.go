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

package lambda

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws/usageReports"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"time"
	"github.com/aws/aws-sdk-go/aws/session"
)

// getFunctionTag formats []*lambda.Tag to map[string]string
func getFunctionTags(ctx context.Context, function *lambda.FunctionConfiguration, svc *lambda.Lambda) []utils.Tag {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	res := make([]utils.Tag, 0)
	tags, err := svc.ListTags(&lambda.ListTagsInput{
		Resource: function.FunctionArn,
	})
	if err != nil {
		logger.Error("Failed to get Lambda tags", err.Error())
		return res
	}
	for key, value := range tags.Tags {
		res = append(res, utils.Tag{
			Key:   key,
			Value: aws.StringValue(value),
		})
	}
	return res
}

// getFunctionInvocationsStats gets the Invocations stats of a Lambda function from CloudWatch
func getFunctionInvocationsStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (Invocations, error) {
	statsTotal, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Invocations"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return Invocations{-1, -1}, err
	}
	statsError, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Errors"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return Invocations{-1, -1}, err
	} else if len(statsTotal.Datapoints) > 0 && len(statsError.Datapoints) > 0 {
		return Invocations{
			aws.Float64Value(statsTotal.Datapoints[0].Sum),
			aws.Float64Value(statsError.Datapoints[0].Sum)}, nil
	} else {
		return Invocations{-1, -1}, nil
	}
}

// getFunctionDurationStats gets the Duration stats of a Lambda function from CloudWatch
func getFunctionDurationStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (Duration, error) {
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Duration"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Average"), aws.String("Maximum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return Duration{-1, -1}, err
	} else if len(stats.Datapoints) > 0 {
		return Duration{
			aws.Float64Value(stats.Datapoints[0].Average),
			aws.Float64Value(stats.Datapoints[0].Maximum)}, nil
	} else {
		return Duration{-1, -1}, nil
	}
}

// getInstanceStats gets the instance stats from CloudWatch
func getFunctionStats(ctx context.Context, instance *lambda.FunctionConfiguration, sess *session.Session, start, end time.Time) (stats Stats) {
	var err error
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := cloudwatch.New(sess)
	dimensions := []*cloudwatch.Dimension{{
		Name:  aws.String("FunctionName"),
		Value: aws.String(aws.StringValue(instance.FunctionName)),
	}}
	stats.Invocations, err = getFunctionInvocationsStats(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching Invocations stats from CloudWatch", err.Error())
	}
	stats.Duration, err = getFunctionDurationStats(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching Duration stats from CloudWatch", err.Error())
	}
	return
}
