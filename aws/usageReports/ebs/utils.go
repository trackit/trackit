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

package ebs

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/es"
)

const MonitorSnapshotStsSessionName = "monitor-snapshot"

type (
	// SnapshotReport is saved in ES to have all the information of an EBS snapshot
	SnapshotReport struct {
		utils.ReportBase
		Snapshot Snapshot `json:"snapshot"`
	}

	// SnapshotBase contains basics information of an EBS snapshot
	SnapshotBase struct {
		Id          string    `json:"id"`
		Description string    `json:"description"`
		State       string    `json:"state"`
		Encrypted   bool      `json:"encrypted"`
		StartTime   time.Time `json:"startTime"`
		Region      string    `json:"region"`
	}

	// Snapshot contains all the information of an EBS snapshot
	Snapshot struct {
		SnapshotBase
		Tags   []utils.Tag `json:"tags"`
		Volume Volume      `json:"volume"`
		Cost   float64     `json:"cost"`
	}

	// Volume contains information about an EBS volume
	Volume struct {
		Id   string `json:"id"`
		Size int64  `json:"size"`
	}
)

// importSnapshotsToEs imports EBS snapshots in ElasticSearch.
// It calls createIndexEs if the index doesn't exist.
func importSnapshotsToEs(ctx context.Context, aa taws.AwsAccount, snapshots []SnapshotReport) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Updating EBS snapshots for AWS account.", map[string]interface{}{
		"awsAccount": aa,
	})
	index := es.IndexNameForUserId(aa.UserId, IndexPrefixEBSReport)
	bp, err := utils.GetBulkProcessor(ctx)
	if err != nil {
		logger.Error("Failed to get bulk processor.", err.Error())
		return err
	}
	for _, snapshot := range snapshots {
		id, err := generateId(snapshot)
		if err != nil {
			logger.Error("Error when marshaling snapshot var", err.Error())
			return err
		}
		bp = utils.AddDocToBulkProcessor(bp, snapshot, TypeEBSReport, index, id)
	}
	bp.Flush()
	err = bp.Close()
	if err != nil {
		logger.Error("Fail to put EBS snapshots in ES", err.Error())
		return err
	}
	logger.Info("EBS snapshots put in ES", nil)
	return nil
}

func generateId(snapshot SnapshotReport) (string, error) {
	ji, err := json.Marshal(struct {
		Account    string    `json:"account"`
		ReportDate time.Time `json:"reportDate"`
		Id         string    `json:"id"`
	}{
		snapshot.Account,
		snapshot.ReportDate,
		snapshot.Snapshot.Id,
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
func merge(cs ...<-chan Snapshot) <-chan Snapshot {
	var wg sync.WaitGroup
	out := make(chan Snapshot)

	// Start an output goroutine for each input channel in cs. The output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Snapshot) {
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
