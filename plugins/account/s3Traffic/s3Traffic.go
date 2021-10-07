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

package plugins_account_s3_traffic

import (
	"fmt"
	"time"

	ts3 "github.com/trackit/trackit/aws/s3"
	"github.com/trackit/trackit/es"
	core "github.com/trackit/trackit/plugins/account/core"
	utils "github.com/trackit/trackit/plugins/utils"
)

func init() {
	// Register the plugin
	core.AccountPlugin{
		Name:            "S3 traffic",
		Description:     "Get the list of s3 buckets with no traffic over the last month",
		Category:        utils.PluginsCategories["S3"],
		Label:           "bucket(s) with traffic",
		Func:            handlerS3Traffic,
		BillingDataOnly: true,
	}.Register()
}

// prepareResult sets the Result and Status in the pluginRes struct
func prepareResult(pluginRes *core.PluginResult) {
	if pluginRes.Checked == pluginRes.Passed {
		pluginRes.Status = "green"
		pluginRes.Result = "All your S3 buckets have traffic"
		return
	}
	pluginRes.Result = fmt.Sprintf("You have %d s3 buckets without traffic", pluginRes.Checked-pluginRes.Passed)
	pluginRes.Status = utils.StatusPercentSteps{50, 80}.GetStatus(pluginRes.Checked, pluginRes.Passed)
}

// getBucketsWithNoTraffic searches for buckets with no traffic and fills the pluginRes struct
func getBucketsWithNoTraffic(pluginRes *core.PluginResult, storage, bandwidth bucketsInfos) {
	for bucketName := range storage {
		pluginRes.Checked++
		if _, ok := bandwidth[bucketName]; ok {
			pluginRes.Passed++
		} else {
			pluginRes.Details = append(pluginRes.Details, bucketName)
		}
	}
	prepareResult(pluginRes)
}

// processS3Traffic retrieves storage and bandwidth information from ES
func processS3Traffic(pluginParams core.PluginParams, pluginRes *core.PluginResult) {
	beginDate := time.Now().AddDate(0, -1, 0).UTC()
	endDate := time.Now().UTC()
	esIndex := es.IndexNameForUserId(pluginParams.User.Id, ts3.IndexPrefixLineItem)

	searchService := GetS3StorageUsage(beginDate, endDate, pluginParams.ESClient, pluginParams.AccountId, esIndex)
	res, err := searchService.Do(pluginParams.Context)
	if err != nil {
		pluginRes.Status = "red"
		pluginRes.Error = fmt.Sprintln("Unable to retrieve S3 storage usage : ", err.Error())
		return
	}
	storage, err := parseESResult(pluginParams, res)
	if err != nil {
		pluginRes.Status = "red"
		pluginRes.Error = fmt.Sprintln("Unable to parse S3 storage usage: ", err.Error())
		return
	}

	searchService = GetS3BandwidthUsage(beginDate, endDate, pluginParams.ESClient, pluginParams.AccountId, esIndex)
	res, err = searchService.Do(pluginParams.Context)
	if err != nil {
		pluginRes.Status = "red"
		pluginRes.Error = fmt.Sprintln("Unable to retrieve S3 bandwidth usage : ", err.Error())
		return
	}
	bandwidth, err := parseESResult(pluginParams, res)
	if err != nil {
		pluginRes.Status = "red"
		pluginRes.Error = fmt.Sprintln("Unable to parse S3 bandwidth usage: ", err.Error())
		return
	}
	getBucketsWithNoTraffic(pluginRes, storage, bandwidth)
}

// handlerS3Traffic is the handler function for the S3 traffic plugin
// it takes a core.PluginParams struct and returns a core.PluginResult struct
func handlerS3Traffic(params core.PluginParams) core.PluginResult {
	res := core.PluginResult{}
	processS3Traffic(params, &res)
	return res
}
