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

package mediastore

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/es"
)

const MonitorContainerStsSessionName = "monitor-instance"

type (
	// ContainerReport is saved in ES to have all the information of an MediaStore instance
	ContainerReport struct {
		utils.ReportBase
		Container Container `json:"container"`
	}

	// ContainerBase contains basics information of an MediaStore instance
	ContainerBase struct {
		Arn    string `json:"arn"`
		Region string `json:"region"`
		Name   string `json:"name"`
	}

	// Container contains all the information of an MediaStore instance
	Container struct {
		ContainerBase
		Costs map[time.Time]float64 `json:"costs"`
	}
)

// importContainersToEs imports MediaStore instances in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importContainersToEs(ctx context.Context, aa taws.AwsAccount, instances []ContainerReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating MediaStore instances for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixMediaStoreReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, instance := range instances {
		id, err := generateId(instance)
		if err != nil {
			logger.Error("Error when marshaling instance var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, instance, TypeMediaStoreReport, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Fail to put MediaStore instances in ES", err.Error())
		return err
	}
	logger.Info("MediaStore instances put in ES", nil)
	return nil
}

func generateId(instance ContainerReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
	}{
		instance.Account,
		instance.ReportDate,
		instance.Container.Name,
	})
	if err != nil {
		return "", err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	return hash64, nil
}

// merge function from https://blog.golang.org/pipelines#TOC_4
// It allows to merge many chans to one.
func merge(cs ...<-chan Container) <-chan Container {
	var wg sync.WaitGroup
	out := make(chan Container)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Container) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done. This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// getPlatformName normalizes the platform name
func getPlatformName(platform string) string {
	if platform == "" {
		return "Linux/UNIX"
	}
	return platform
}
