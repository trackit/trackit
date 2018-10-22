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

package ec2

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/trackit/jsonlog"
)

// getPurchasingOption returns a string that describes how the instance given as parameter have been purchased
func getPurchasingOption(instance *ec2.Instance) string {
	var purchasing string
	lifeCycle := aws.StringValue(instance.InstanceLifecycle)
	tenancy := aws.StringValue(instance.Placement.Tenancy)
	if tenancy == "" || tenancy == "default" {
		if lifeCycle == "" {
			purchasing = "on demand"
		} else {
			purchasing = lifeCycle
		}
	} else {
		purchasing = tenancy
	}
	return purchasing
}

// getInstanceTag formats []*ec2.Tag to map[string]string
func getInstanceTag(tags []*ec2.Tag) map[string]string {
	res := make(map[string]string)
	for _, tag := range tags {
		res[aws.StringValue(tag.Key)] = aws.StringValue(tag.Value)
	}
	return res
}

// getInstanceCPUStats gets the CPU average and the CPU peak of an EC2 instance from CloudWatch
func getInstanceCPUStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, error) {
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("CPUUtilization"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
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

// getInstanceNetworkStats gets the network in and out stats of an EC2 instance from CloudWatch
func getInstanceNetworkStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, error) {
	statsIn, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("NetworkIn"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return -1, -1, err
	}
	statsOut, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("NetworkOut"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return -1, -1, err
	} else if len(statsIn.Datapoints) > 0 && len(statsOut.Datapoints) > 0 {
		return aws.Float64Value(statsIn.Datapoints[0].Sum), aws.Float64Value(statsOut.Datapoints[0].Sum), nil
	} else {
		return -1, -1, nil
	}
}

// getInstanceInternalIOStats gets the IO read and write stats of an EC2 instance from CloudWatch
func getInstanceInternalIOStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, error) {
	statsRead, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("DiskReadBytes"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return -1, -1, err
	}
	statsWrite, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("DiskWriteBytes"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return -1, -1, err
	} else if len(statsRead.Datapoints) > 0 && len(statsWrite.Datapoints) > 0 {
		return aws.Float64Value(statsRead.Datapoints[0].Sum), aws.Float64Value(statsWrite.Datapoints[0].Sum), nil
	} else {
		return -1, -1, nil
	}
}

// getInstanceELBIOStats gets the IO read and write stats of a volume from CloudWatch
func getInstanceELBIOStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, error) {
	statsRead, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EBS"),
		MetricName: aws.String("VolumeReadBytes"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return -1, -1, err
	}
	statsWrite, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EBS"),
		MetricName: aws.String("VolumeWriteBytes"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Period:     aws.Int64(int64(60*60*24) * 31),
		Statistics: []*string{aws.String("Sum")},
		Dimensions: dimensions,
	})
	if err != nil {
		return -1, -1, err
	} else if len(statsRead.Datapoints) > 0 && len(statsWrite.Datapoints) > 0 {
		return aws.Float64Value(statsRead.Datapoints[0].Sum), aws.Float64Value(statsWrite.Datapoints[0].Sum), nil
	} else {
		return -1, -1, nil
	}
}

// getInstanceIOStats gets the IO read and write stats from CloudWatch
func getInstanceIOStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension,
	volumes []string, start, end time.Time) (map[string]float64, map[string]float64, error) {
	statsRead := make(map[string]float64, 0)
	statsWrite := make(map[string]float64, 0)
	internalRead, internalWrite, err := getInstanceInternalIOStats(svc, dimensions, start, end)
	if err != nil {
		return nil, nil, err
	}
	statsRead["internal"] = float64(internalRead)
	statsWrite["internal"] = float64(internalWrite)
	for _, volume := range volumes {
		dimensionsEBS := []*cloudwatch.Dimension{{
			Name:  aws.String("VolumeId"),
			Value: aws.String(volume),
		},
		}
		read, write, err := getInstanceELBIOStats(svc, dimensionsEBS, start, end)
		if err != nil {
			return nil, nil, err
		}
		statsRead[string(volume)] = float64(read)
		statsWrite[string(volume)] = float64(write)
	}
	return statsRead, statsWrite, nil
}

// getInstanceStats gets the instance stats from CloudWatch
func getInstanceStats(ctx context.Context, instance *ec2.Instance, sess *session.Session, start, end time.Time) instanceStats {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := cloudwatch.New(sess)
	dimensions := []*cloudwatch.Dimension{{
		Name:  aws.String("InstanceId"),
		Value: aws.String(aws.StringValue(instance.InstanceId)),
	},
	}
	CpuAverage, CpuPeak, err := getInstanceCPUStats(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching CPU stats from CloudWatch", err.Error())
	}
	NetworkIn, NetworkOut, err := getInstanceNetworkStats(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching Network stats from CloudWatch", err.Error())
	}
	volumes := make([]string, 0)
	for _, volume := range instance.BlockDeviceMappings {
		volumes = append(volumes, aws.StringValue(volume.Ebs.VolumeId))
	}
	IORead, IOWrite, err := getInstanceIOStats(svc, dimensions, volumes, start, end)
	if err != nil {
		logger.Error("Error when fetching IO stats from CloudWatch", err.Error())
	}
	return instanceStats{
		CpuAverage,
		CpuPeak,
		NetworkIn,
		NetworkOut,
		IORead,
		IOWrite,
	}
}
