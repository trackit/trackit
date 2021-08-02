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

package plugins_account_network_ec2

import (
	"fmt"
	"time"

	"github.com/trackit/trackit/db"
	core "github.com/trackit/trackit/plugins/account/core"
	utils "github.com/trackit/trackit/plugins/utils"
	"github.com/trackit/trackit/usageReports/ec2"
)

const (
	// 10 Mebibyte
	networkLimit = 1.049e+7
)

func init() {
	// Register the plugin
	core.AccountPlugin{
		Name:        "EC2 Network",
		Description: "Get the list of EC2 instances with low network activity",
		Category:    utils.PluginsCategories["EC2"],
		Label:       "EC2 instance(s) with network activity",
		Func:        processNetworkEc2,
	}.Register()
}

// prepareResult sets the Result and Status in the pluginRes struct
func prepareResult(pluginRes *core.PluginResult) {
	if pluginRes.Checked == pluginRes.Passed {
		pluginRes.Status = "green"
		pluginRes.Result = "All your EC2 instances have network activity"
	} else {
		pluginRes.Result = fmt.Sprintf("You have %d EC2 instance with low network activity", pluginRes.Checked-pluginRes.Passed)
		pluginRes.Status = utils.StatusPercentSteps{50, 80}.GetStatus(pluginRes.Checked, pluginRes.Passed)
	}
}

// getUnusedEc2Recommendation searches for unused ec2 network in usage report
// It takes a *core.PluginResult and an array of instances as parameters
func getUnusedEc2Recommendation(pluginRes *core.PluginResult, instances []ec2.InstanceReport) {
	pluginRes.Details = make([]string, 0)
	for _, instance := range instances {
		if instance.Instance.Stats.Network.In == -1 || instance.Instance.Stats.Network.Out == -1 {
			continue
		}
		pluginRes.Checked++
		network := instance.Instance.Stats.Network.In + instance.Instance.Stats.Network.Out
		if network > networkLimit {
			pluginRes.Passed++
		} else {
			pluginRes.Details = append(pluginRes.Details, fmt.Sprintf("%s %s", instance.Instance.Id, instance.Instance.Tags["Name"]))
		}
	}
	prepareResult(pluginRes)
}

// processNetworkEc2 is the handler function for the Unused EC2 Network plugin
// it takes a core.PluginParams struct and returns a core.PluginResult struct
func processNetworkEc2(params core.PluginParams) (res core.PluginResult) {
	tx, err := db.Db.Begin()
	if err != nil {
		res.Status = "red"
		res.Error = fmt.Sprintln("Unable to retrieve EC2 instances : %s", err.Error())
		return
	}
	_, instances, err := ec2.GetEc2Data(params.Context,
		ec2.Ec2QueryParams{[]string{params.AccountId}, nil, time.Now().UTC()},
		params.User, tx)
	if err != nil {
		res.Status = "red"
		res.Error = fmt.Sprintln("Unable to retrieve EC2 instances : %s", err.Error())
		return
	}
	getUnusedEc2Recommendation(&res, instances)
	return
}
