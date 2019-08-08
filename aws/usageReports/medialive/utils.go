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

package medialive

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

const MonitorChannelStsSessionName = "monitor-channel"

type (
	// ChannelReport is saved in ES to have all the information of an MediaLive channel
	ChannelReport struct {
		utils.ReportBase
		Channel Channel `json:"channel"`
	}

	// ChannelBase contains basics information of an MediaLive channel
	ChannelBase struct {
		Arn    string `json:"arn"`
		Id     string `json:"id"`
		Name   string `json:"name"`
		Region string `json:"region"`
	}

	// Channel contains all the information of an MediaLive channel
	Channel struct {
		ChannelBase
		Cost map[time.Time]float64 `json:"cost"`
		Tags map[string]string     `json:"tags"`
	}
	InputReport struct {
		utils.ReportBase
		Input Input `json:"channel"`
	}

	// InputBase contains basics information of an MediaLive channel
	InputBase struct {
		Arn    string `json:"arn"`
		Id     string `json:"id"`
		Name   string `json:"name"`
		Region string `json:"region"`
	}

	// Input contains all the information of an MediaLive channel
	Input struct {
		InputBase
		Cost map[time.Time]float64 `json:"cost"`
		Tags map[string]string     `json:"tags"`
	}
)

// importChannelsToEs imports MediaLive channels in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importChannelsToEs(ctx context.Context, aa taws.AwsAccount, channels []ChannelReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating MediaLive channels for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixMediaLiveReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, channel := range channels {
		id, err := generateId(channel)
		if err != nil {
			logger.Error("Error when marshaling channel var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, channel, TypeMediaLiveReport, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Fail to put MediaLive channels in ES", err.Error())
		return err
	}
	logger.Info("MediaLive channels put in ES", nil)
	return nil
}

// importChannelsToEs imports MediaLive channels in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importInputsToEs(ctx context.Context, aa taws.AwsAccount, inputs []InputReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating MediaLive inputs for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixMediaLiveInputReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, input := range inputs {
		id, err := generateInputId(input)
		if err != nil {
			logger.Error("Error when marshaling input var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, input, TypeMediaLiveInputReport, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Fail to put MediaLive inputs in ES", err.Error())
		return err
	}
	logger.Info("MediaLive inputs put in ES", nil)
	return nil
}

func generateId(channel ChannelReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
	}{
		channel.Account,
		channel.ReportDate,
		channel.Channel.Id,
	})
	if err != nil {
		return "", err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	return hash64, nil
}

func generateInputId(input InputReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
	}{
		input.Account,
		input.ReportDate,
		input.Input.Id,
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
func mergeChannels(cs ...<-chan Channel) <-chan Channel {
	var wg sync.WaitGroup
	out := make(chan Channel)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Channel) {
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

func mergeInput(cs ...<-chan Input) <-chan Input {
	var wg sync.WaitGroup
	out := make(chan Input)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Input) {
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
