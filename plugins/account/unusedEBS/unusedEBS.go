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

package plugins_account_unused_ebs

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/trackit/trackit-server/config"
	core "github.com/trackit/trackit-server/plugins/account/core"
	"github.com/trackit/trackit-server/plugins/utils"
)

func init() {
	// Register the plugin
	core.AccountPlugin{
		Name:        "unused_ebs",
		Description: "unused ebs plugin",
		Func:        processUnusedEBS,
	}.Register()
}

// formatResult takes a map of unused EBS and returns 3 strings to fill the
// core.PluginResult struct
func formatResult(unusedByAZ map[string]int) (string, string, string) {
	total := 0
	detail := new(bytes.Buffer)
	for az, totalAz := range unusedByAZ {
		if total == 0 {
			fmt.Fprintf(detail, fmt.Sprintf("%s(%d)", az, totalAz))
		} else {
			fmt.Fprintf(detail, fmt.Sprintf(", %s(%d)", az, totalAz))
		}
		total += totalAz
	}
	if total == 0 {
		return "You don't have any unused EBS", "", ""
	}
	return fmt.Sprintf("You have %d unused EBS", total), detail.String(), ""
}

// getUnusedEBsRecommendation searches for unused ebs in every region available
// It takes a core.PluginParams struct and returns 3 strings (result, details, error)
func getUnusedEBsRecommendation(pluginParams core.PluginParams) (string, string, string) {
	svc := plugins_utils.GetEc2ClientSession(pluginParams.AccountCredentials, &config.AwsRegion)
	regionsOutput, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return "", "", fmt.Sprintf("Unable to retrieve the list of regions: %s", err.Error())
	}

	unusedByAZ := make(map[string]int)
	for _, region := range regionsOutput.Regions {
		svc = plugins_utils.GetEc2ClientSession(pluginParams.AccountCredentials, region.RegionName)
		err = svc.DescribeVolumesPages(&ec2.DescribeVolumesInput{},
			func(page *ec2.DescribeVolumesOutput, lastPage bool) bool {
				for _, volume := range page.Volumes {
					if volume != nil && *volume.State == "available" {
						unusedByAZ[*volume.AvailabilityZone] = unusedByAZ[*volume.AvailabilityZone] + 1
					}
				}
				return lastPage
			})
		if err != nil {
			return "", "", fmt.Sprintf("Unable to list volumes: %s", err.Error())
		}
	}
	return formatResult(unusedByAZ)
}

// processUnusedEBS is the handler function for the unused_ebs plugin
// it takes a core.PluginParams struct and returns a core.PluginResult struct
func processUnusedEBS(params core.PluginParams) core.PluginResult {
	result, details, err := getUnusedEBsRecommendation(params)
	return core.PluginResult{
		Result:  result,
		Details: details,
		Error:   err,
	}
}
