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

package plugins_account_anattached_eip

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/trackit/trackit/config"
	core "github.com/trackit/trackit/plugins/account/core"
	utils "github.com/trackit/trackit/plugins/utils"
)

func init() {
	// Register the plugin
	core.AccountPlugin{
		Name:        "Unattached EIP",
		Description: "Returns the list of unattached EIP",
		Category:    utils.PluginsCategories["EC2"],
		Label:       "attached EIP(s)",
		Func:        processUnattachedEIP,
	}.Register()
}

// prepareResult sets the Result and Status in the pluginRes struct
func prepareResult(pluginRes *core.PluginResult) {
	if pluginRes.Checked == pluginRes.Passed {
		pluginRes.Status = "green"
		pluginRes.Result = "You don't have any unattached EIP"
		return
	}
	pluginRes.Result = fmt.Sprintf("You have %d unattached EIP", pluginRes.Checked-pluginRes.Passed)
	pluginRes.Status = utils.StatusPercentSteps{50, 95}.GetStatus(pluginRes.Checked, pluginRes.Passed)
}

// processEIP checks if the EIP for a given region are attached and fills the pluginRes struct accordingly
func processEIP(pluginRes *core.PluginResult, region *string, eipRes *ec2.DescribeAddressesOutput) {
	if eipRes.Addresses != nil {
		for _, eip := range eipRes.Addresses {
			pluginRes.Checked += 1
			if eip.AssociationId != nil {
				pluginRes.Passed += 1
			} else {
				eipDesc := aws.StringValue(eip.PublicIp)
				if eipDesc == "" {
					eipDesc = aws.StringValue(eip.AssociationId)
				}
				pluginRes.Details = append(pluginRes.Details, fmt.Sprintf("%s (%s)", eipDesc, *region))
			}
		}
	}
}

// fetchEIPInfos fetches EIP address infos for a given region
func fetchEIPInfos(pluginParams core.PluginParams, region *string, eipChan chan EIP) error {
	defer close(eipChan)
	svc := utils.GetEc2ClientSession(pluginParams.AccountCredentials, region)
	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{})
	eip := EIP{
		Region: region,
	}
	if err != nil {
		eip.Err = err
		eipChan <- eip
		return err
	}
	eip.EIPRes = result
	eipChan <- eip
	return nil
}

// getUnattachedEIP searches for unused EIP in every region available
// It takes a core.PluginParams struct and a *core.PluginResult as parameters
func getUnattachedEIP(pluginParams core.PluginParams, pluginRes *core.PluginResult) {
	svc := utils.GetEc2ClientSession(pluginParams.AccountCredentials, &config.AwsRegion)
	regionsOutput, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		pluginRes.Status = "red"
		pluginRes.Error = fmt.Sprintf("Unable to retrieve the list of regions: %s", err.Error())
		return
	}
	EIPChans := make([]<-chan EIP, 0, len(regionsOutput.Regions))
	for _, region := range regionsOutput.Regions {
		eipChan := make(chan EIP)
		go fetchEIPInfos(pluginParams, region.RegionName, eipChan) // No need to handle the error here, it's sent down eipChan and thus handled in the next loop
		EIPChans = append(EIPChans, eipChan)
	}
	for eip := range merge(EIPChans...) {
		if eip.Err != nil {
			pluginRes.Status = "red"
			pluginRes.Error = fmt.Sprintf("Unable to list addresses: %s", eip.Err.Error())
			return
		}
		processEIP(pluginRes, eip.Region, eip.EIPRes)
	}
	prepareResult(pluginRes)
}

// processUnattachedEIP is the handler function for the Unattached EIP plugin
// it takes a core.PluginParams struct and returns a core.PluginResult struct
func processUnattachedEIP(params core.PluginParams) core.PluginResult {
	res := core.PluginResult{}
	getUnattachedEIP(params, &res)
	return res
}
