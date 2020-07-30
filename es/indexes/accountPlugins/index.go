//   Copyright 2020 MSolution.IO
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

package accountPlugins

import "time"

const IndexSuffix = "account-plugins"
const Type = "account-plugin"
const TemplateName = "account-plugins"

// PluginResultES is the struct used to save a plugin result into elaticsearch
type PluginResultES struct {
	AccountPluginIdx string    `json:"accountPluginIdx"`
	Account          string    `json:"account"`
	ReportDate       time.Time `json:"reportDate"`
	PluginName       string    `json:"pluginName"`
	Category         string    `json:"category"`
	Label            string    `json:"label"`
	Result           string    `json:"result"`
	Status           string    `json:"status"`
	Details          []string  `json:"details"`
	Error            string    `json:"error"`
	Checked          int       `json:"checked"`
	Passed           int       `json:"passed"`
}
