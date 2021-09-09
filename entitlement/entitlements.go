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

package entitlement

import (
	"context"
	"database/sql"

	"github.com/trackit/jsonlog"
)

type entitlementConfig struct {
	Get    getUserEntitlementFunc
	Update updateUserEntitlementFunc
	Type   string
}

var entitlementConfigs = []entitlementConfig{
	{
		getUserEntitlementMarketplace,
		updateUserEntitlementMarketplace,
		"Marketplace",
	},
	{
		getUserEntitlementTagbotMarketplace,
		updateUserEntitlementTagbotMarketplace,
		"Marketplace for Tagbot",
	},
}

func CheckUserEntitlements(ctx context.Context, db *sql.Tx, userId int) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	for _, entitlementConf := range entitlementConfigs {
		logger.Info("Checking user entitlement", map[string]interface{}{
			"entitlement": entitlementConf.Type,
			"userId":      userId,
		})
		err := checkUserEntitlement(ctx, db, userId, entitlementConf.Get, entitlementConf.Update)
		if err != nil {
			logger.Error("Could not check user entitlement", map[string]interface{}{
				"entitlement": entitlementConf.Type,
				"err":         err.Error(),
			})
		}
	}

	return nil
}

type getUserEntitlementFunc func(db *sql.Tx, ctx context.Context, userId int) (bool, error)
type updateUserEntitlementFunc func(db *sql.Tx, ctx context.Context, userId int, entitlementValue bool) error

func checkUserEntitlement(ctx context.Context, db *sql.Tx, userId int, get getUserEntitlementFunc, update updateUserEntitlementFunc) error {
	entitlement, err := get(db, ctx, userId)
	if err != nil {
		return err
	}
	return update(db, ctx, userId, entitlement)
}
