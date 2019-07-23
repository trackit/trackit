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

package elasticache

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws/usageReports"
)

func getClusterTags(ctx context.Context, cluster *elasticache.CacheCluster, svc *elasticache.ElastiCache, account, region string) []utils.Tag {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	tags := make([]utils.Tag, 0)
	// format ARN for an ElastiCache instance is:
	// "arn:aws:elasticache:[region]:[aws_id]:cluster:[cluster name]"
	arn := "arn:aws:elasticache:" + region + ":" + account + ":cluster:" + aws.StringValue(cluster.CacheClusterId)
	res, err := svc.ListTagsForResource(&elasticache.ListTagsForResourceInput{
		ResourceName: aws.String(arn),
	})
	if err != nil {
		logger.Error("Error while getting cluster tags", err.Error())
		return tags
	}
	for _, tag := range res.TagList {
		tags = append(tags, utils.Tag{
			Key:   aws.StringValue(tag.Key),
			Value: aws.StringValue(tag.Value),
		})
	}
	return tags
}

func extractCacheNodes(cacheNodes []*elasticache.CacheNode) []Node {
	nodes := make([]Node, 0)
	for _, node := range cacheNodes {
		nodes = append(nodes, Node{
			Id:     aws.StringValue(node.CacheNodeId),
			Status: aws.StringValue(node.CacheNodeStatus),
			Region: aws.StringValue(node.CustomerAvailabilityZone),
		})
	}
	return nodes
}

// getInstanceCPUStats gets the CPU average and the CPU peak of an ElastiCache instance from CloudWatch
func getInstanceCPUStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, error) {
	stats, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ElastiCache"),
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

// getInstanceNetworkStats gets the network in and out stats of an ElastiCache instance from CloudWatch
func getInstanceNetworkStats(svc *cloudwatch.CloudWatch, dimensions []*cloudwatch.Dimension, start, end time.Time) (float64, float64, error) {
	statsIn, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ElastiCache"),
		MetricName: aws.String("NetworkBytesIn"),
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
		Namespace:  aws.String("AWS/ElastiCache"),
		MetricName: aws.String("NetworkBytesOut"),
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

// getInstanceStats gets the instance stats from CloudWatch
func getInstanceStats(ctx context.Context, instance *elasticache.CacheCluster, sess *session.Session, start, end time.Time) Stats {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	svc := cloudwatch.New(sess)
	dimensions := []*cloudwatch.Dimension{{
		Name:  aws.String("CacheClusterId"),
		Value: aws.String(aws.StringValue(instance.CacheClusterId)),
	}}
	var stats Stats
	var err error = nil
	stats.Cpu.Average, stats.Cpu.Peak, err = getInstanceCPUStats(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching CPU stats from CloudWatch", err.Error())
	}
	stats.Network.In, stats.Network.Out, err = getInstanceNetworkStats(svc, dimensions, start, end)
	if err != nil {
		logger.Error("Error when fetching Network stats from CloudWatch", err.Error())
	}
	return stats
}
