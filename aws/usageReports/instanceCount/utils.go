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

package instanceCount

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/aws/usageReports"
)

const MonitorInstanceCountStsSessionName = "monitor-instanceCount"

type (
	// InstanceCount is saved in ES to have all the information of an InstanceCount
	InstanceCountReport struct {
		utils.ReportBase
		InstanceCount InstanceCount `json:"instanceCount"`
	}

	// InstanceCount contains all the information of an InstanceCount
	InstanceCount struct {
		Type   string              `json:"instanceType"`
		Hours []InstanceCountHours `json:"hours"`
	}

	InstanceCountHours struct {
		Hour  time.Time `json:"hour"`
		Count int       `json:"count"`
	}
)

// importInstanceCountsToEs imports Instance Count in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importInstanceCountToEs(ctx context.Context, aa taws.AwsAccount, reports []InstanceCountReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating InstanceCount for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	//index := es.IndexNameForUserId(aa.UserId, IndexPrefixInstanceCountReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	/*for _, report := range reports {
		if err != nil {
			logger.Error("Error when marshaling instanceCount var", err.Error())
			return err
		}
		//bp = utils.AddDocToBulkProcessor(bp, report, TypeInstanceCountReport, index, id)
	}*/
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Fail to put InstanceCount in ES", err.Error())
		return err
	}
	logger.Info("InstanceCount put in ES", nil)
	return nil
}

func generateId(report InstanceCount) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
	}{
		/*report.Account,
		report.ReportDate,
		report.InstanceCount.Id,*/
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
func merge(cs ...<-chan InstanceCount) <-chan InstanceCount {
	var wg sync.WaitGroup
	out := make(chan InstanceCount)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan InstanceCount) {
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
