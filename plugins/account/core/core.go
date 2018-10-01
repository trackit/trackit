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

package plugins_account_core

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/users"
	"gopkg.in/olivere/elastic.v5"
)

// AccountPlugin is the struct that defines the variables and functions a plugin
// needs to use
type AccountPlugin struct {
	Name        string
	Description string
	Func        PluginFunc
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

// PluginParams is the struct that each plugin should return
type PluginResult struct {
	Result  string
	Details string
	Error   string
}

// PluginResultES is the struct used to save a plugin result into elaticsearch
type PluginResultES struct {
	ReportDate time.Time `json:"reportDate"`
	PluginName string    `json:"pluginName"`
	Account    string    `json:"account"`
	Result     string    `json:"result"`
	Details    string    `json:"details"`
	Error      string    `json:"error"`
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
