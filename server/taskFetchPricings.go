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

package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws/pricings"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/models"
)

// taskFetchPricings fetches the EC2 pricings and saves them in the database
func taskFetchPricings(ctx context.Context) (err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	res, err := pricings.FetchEc2Pricings(ctx)
	if err != nil {
		logger.Error("Failed to retrieve ec2 pricings", err.Error())
		return
	}
	serializedPricing, err := json.Marshal(res)
	if err != nil {
		logger.Error("Failed to serialize ec2 pricings", err.Error())
		return
	}
	var tx *sql.Tx
	defer func() {
		if tx != nil {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
		logger.Error("Failed to initiate sql transaction", err.Error())
		return
	} else {
		ec2PricingDb, _ := models.AwsPricingByProduct(tx, pricings.EC2ServiceCode)
		if ec2PricingDb == nil {
			ec2PricingDb = &models.AwsPricing{
				Product: pricings.EC2ServiceCode,
			}
		}
		ec2PricingDb.Pricing = serializedPricing
		err = ec2PricingDb.Save(tx)
		if err != nil {
			logger.Error("Failed to save ec2 pricings", err.Error())
			return
		}
	}
	return
}
