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
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/accountPlugins"
)

// IngestPluginResult saves a PluginResultES into elasticsearch
func IngestPluginResult(ctx context.Context, aa aws.AwsAccount, pluginRes accountPlugins.PluginResultES) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Saving plugin result for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	client := es.Client
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
	}{
		pluginRes.Account,
		pluginRes.ReportDate,
	})
	if err != nil {
		logger.Error("Error when marshaling instance var", err.Error())
		return err
	}
	hash := md5.Sum(ji)
	hash64 := base64.URLEncoding.EncodeToString(hash[:])
	index := es.IndexNameForUserId(aa.UserId, accountPlugins.IndexSuffix)
	pluginRes.AccountPluginIdx = fmt.Sprintf("%s-%s", pluginRes.Account, pluginRes.PluginName)
	if res, err := client.
		Index().
		Index(index).
		Type(accountPlugins.Type).
		BodyJson(pluginRes).
		Id(hash64).
		Do(context.Background()); err != nil {
		logger.Error("Error when putting plugin result in ES", err.Error())
		return err
	} else {
		logger.Info("plugin result put in ES", *res)
	}
	return nil
}
