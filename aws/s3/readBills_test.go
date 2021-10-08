//   Copyright 2017 MSolution.IO
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

package s3

import (
	"context"
	"testing"

	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
)

func init() {
	jsonlog.DefaultLogger = jsonlog.DefaultLogger.WithLogLevel(jsonlog.LogLevelDebug)
}

/*
func TestEpitechio(t *testing.T) {
	var count int
	var start time.Time
	var end time.Time
	time.Sleep(10 * time.Second)
	err := ReadBills(
		context.Background(),
		// taws.AwsAccount{
		// 	RoleArn:  "arn:aws:iam::895365654851:role/trackit",
		// 	External: "RLuxJFYhaZYjWHNYY_pfeAgF@lzymhUKNxiwq_IQ",
		// },
		taws.AwsAccount{
			RoleArn:  "arn:aws:iam::394125495069:role/delegation-trackit",
			External: "5e22M7zqmMNo7t4G7SUSD9IK",
		},
		// BillRepository{
		// 	Bucket: "epitechio-reports",
		// 	Prefix: "constandusage/",
		// },
		BillRepository{
			Bucket: "trackit-billing-report",
			Prefix: "usagecost/AllHourlyToS3/",
		},
		func(li LineItem, ok bool) bool {
			if ok {
				if count == 0 {
					start = time.Now()
				}
				count++
			} else {
				end = time.Now()
			}
			return true
		},
		acceptAllManifests,
	)
	fmt.Printf("Parsed %d records in %s.", count, end.Sub(start).String())
	if err != nil {
		println(err.Error())
	}
}
*/

/*
func TestEpitechio(t *testing.T) {
	var count int
	var start time.Time
	var end time.Time
	client, err := elastic.NewClient(
		elastic.SetBasicAuth("elastic", "changeme"),
	)
	if err != nil {
		println(err.Error())
		return
	}
	err = ReadBills(
		context.Background(),
		// taws.AwsAccount{
		// 	RoleArn:  "arn:aws:iam::895365654851:role/trackit",
		// 	External: "RLuxJFYhaZYjWHNYY_pfeAgF@lzymhUKNxiwq_IQ",
		// },
		taws.AwsAccount{
			RoleArn:  "arn:aws:iam::394125495069:role/delegation-trackit",
			External: "5e22M7zqmMNo7t4G7SUSD9IK",
		},
		// BillRepository{
		// 	Bucket: "epitechio-reports",
		// 	Prefix: "constandusage/",
		// },
		BillRepository{
			Bucket: "trackit-billing-report",
			Prefix: "usagecost/AllHourlyToS3/",
		},
		addToElasticsearch(client),
		acceptAllManifests,
	)
	fmt.Printf("Parsed %d records in %s.", count, end.Sub(start).String())
	if err != nil {
		println(err.Error())
	}
}
*/

// Seems to just hang forever rn ? Need more investgation on whether this is some local ES problem, or perhaps due to the hardcoded role given here
func TestUpdate(t *testing.T) {
	latestManifest, err := UpdateReport(
		context.Background(),
		taws.AwsAccount{
			RoleArn:  "arn:aws:iam::394125495069:role/delegation-trackit",
			External: "5e22M7zqmMNo7t4G7SUSD9IK",
		},
		BillRepository{
			Bucket: "trackit-billing-report",
			Prefix: "usagecost/AllHourlyToS3/",
		},
	)
	println(latestManifest.String())
	if err != nil {
		println(err.Error())
	}
}
