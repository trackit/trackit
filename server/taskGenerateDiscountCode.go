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

package main

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/models"
)

// makeRandom10CharString generates a 10-character random string composed of alphanumeric characters
func makeRandom10CharString() string {
	const lettersBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, 10)
	for i := range result {
		result[i] = lettersBytes[rnd.Intn(len(lettersBytes))]
	}
	return string(result)
}

// taskGenerateDiscountCode generates a new random discount code with the given description
func taskGenerateDiscountCode(ctx context.Context) (err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	args := paramsFromContextOrArgs(ctx)
	if len(args) != 1 {
		return errors.New("taskGenerateDiscountCode requires 1 argument")
	}
	description := args[0]

	newDiscountCode := models.TagbotDiscountCode{
		Code:        makeRandom10CharString(),
		Description: description,
	}

	var tx *sql.Tx
	defer utilsUsualTxFinalize(&tx, &err, &logger, "generate-discount-code")
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
		logger.Error("Failed to initiate sql transaction", err.Error())
		return
	}
	if err = newDiscountCode.Insert(tx); err != nil {
		logger.Error("Failed to create new discount code", err.Error())
		return
	}

	jsonlog.LoggerFromContextOrDefault(ctx).Info("Created new discount code", map[string]interface{}{
		"discountCode": newDiscountCode.Code,
	})
	return
}
