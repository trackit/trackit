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

package history

import (
	"time"
	"context"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
)

func getHistoryDate() (time.Time, time.Time) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month() - 1, 0, 0, 0, 0, 0, now.Location())
	return start, end
}

func FetchInfos(ctx context.Context, aa aws.AwsAccount) (error, error) {
	startDate, endDate := getHistoryDate()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Fetching EC2 and RDS history for " + string(aa.Id) + " (" + aa.Pretty + ")", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate,
		"endDate":      endDate,
	})
	return nil, nil
}