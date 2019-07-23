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

package plugins_account_unused_ebs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/trackit/trackit/config"
	core "github.com/trackit/trackit/plugins/account/core"
	utils "github.com/trackit/trackit/plugins/utils"
)

func init() {
	// Register the plugin
	core.AccountPlugin{
		Name:        "Unused EBS",
		Description: "unused ebs plugin",
		Category:    utils.PluginsCategories["EC2"],
		Label:       "attached EBS volume(s)",
		Func:        processUnusedEBS,
	}.Register()
}

// prepareResult takes a map of unused EBS and a *core.PluginResult as parameters
// and fills the PluginResult
func prepareResult(unusedByAZ map[string]int, pluginRes *core.PluginResult) {
	total := 0
	for az, totalAz := range unusedByAZ {
		pluginRes.Details = append(pluginRes.Details, fmt.Sprintf("%s: %d unused volume(s)", az, totalAz))
		total += totalAz
	}
	if total == 0 {
		pluginRes.Result = "You don't have any unused EBS"
		pluginRes.Status = "green"
		return
	}
	pluginRes.Result = fmt.Sprintf("You have %d unused EBS", total)
	pluginRes.Status = utils.StatusPercentSteps{50, 95}.GetStatus(pluginRes.Checked, pluginRes.Passed)
}

// getUnusedEBsRecommendation searches for unused ebs in every region available
// It takes a core.PluginParams struct and a *core.PluginResult as parameters
func getUnusedEBsRecommendation(pluginParams core.PluginParams, pluginRes *core.PluginResult) {
	svc := utils.GetEc2ClientSession(pluginParams.AccountCredentials, &config.AwsRegion)
	regionsOutput, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		pluginRes.Status = "red"
		pluginRes.Error = fmt.Sprintf("Unable to retrieve the list of regions: %s", err.Error())
		return
	}

	unusedByAZ := make(map[string]int)
	for _, region := range regionsOutput.Regions {
		svc = utils.GetEc2ClientSession(pluginParams.AccountCredentials, region.RegionName)
		err = svc.DescribeVolumesPages(&ec2.DescribeVolumesInput{},
			func(page *ec2.DescribeVolumesOutput, lastPage bool) bool {
				for _, volume := range page.Volumes {
					pluginRes.Checked += 1
					if volume != nil && *volume.State == "available" {
						unusedByAZ[*volume.AvailabilityZone] = unusedByAZ[*volume.AvailabilityZone] + 1
					} else {
						pluginRes.Passed += 1
					}
				}
				return lastPage
			})
		if err != nil {
			pluginRes.Status = "red"
			pluginRes.Error = fmt.Sprintf("Unable to list volumes: %s", err.Error())
			return
		}
	}
	prepareResult(unusedByAZ, pluginRes)
}

// processUnusedEBS is the handler function for the Unused EBS plugin
// it takes a core.PluginParams struct and returns a core.PluginResult struct
func processUnusedEBS(params core.PluginParams) core.PluginResult {
	res := core.PluginResult{}
	getUnusedEBsRecommendation(params, &res)
	return res
}
