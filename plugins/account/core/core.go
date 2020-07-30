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

package plugins_account_core

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/olivere/elastic"
	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/users"
)

// AccountPlugin is the struct that defines the variables and functions a plugin
// needs to use
type AccountPlugin struct {
	Name            string
	Description     string
	Category        string
	Label           string
	Func            PluginFunc
	BillingDataOnly bool
}

// PluginParams is the struct that is passed as a parameter for each plugin
type PluginParams struct {
	Context            context.Context
	User               users.User
	AwsAccount         aws.AwsAccount
	AccountId          string
	AccountCredentials *credentials.Credentials
	ESClient           *elastic.Client
}

// PluginResult is the struct that each plugin should return
type PluginResult struct {
	Result  string
	Status  string
	Details []string
	Error   string
	Checked int
	Passed  int
}

// PluginFunc is the type that should be implemented by the plugin's function
type PluginFunc func(PluginParams) PluginResult

// RegisteredAccountPlugins is the list of registered plugins
var RegisteredAccountPlugins = make([]AccountPlugin, 0, 0x40)

// Register allows plugins to register themselves on server startup
func (ap AccountPlugin) Register() AccountPlugin {
	RegisteredAccountPlugins = append(RegisteredAccountPlugins, ap)
	return ap
}
