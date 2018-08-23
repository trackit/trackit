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
	"fmt"
	"time"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
	"github.com/trackit/trackit-server/costs/anomalies"
	"github.com/trackit/trackit-server/db"
	"github.com/trackit/trackit-server/mail"
	"github.com/trackit/trackit-server/models"
	"github.com/trackit/trackit-server/users"
)

// taskAnomaliesDetection processes an AwsAccount to email
// the user if anomalies are detected.
func taskAnomaliesDetection(ctx context.Context) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Debug("Running task 'anomalies-detection'.", nil)
	return EmailAnomalies(ctx)
}

// sendAnomalyEmail actually sends the mail and log everything.
func sendAnomalyEmail(user users.User, awsAccount aws.AwsAccount, product string, an anomalies.CostAnomaly, date time.Time, ctx context.Context) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	ea := models.EmailedAnomaly{
		Date:         date,
		Recipient:    user.Email,
		AwsAccountID: awsAccount.Id,
		Product:      product,
	}
	subject := fmt.Sprintf("TrackIt detected an abnormal peak in your %s cost on your account %s!", product, awsAccount.Pretty)
	body := fmt.Sprintf(`
TrackIt detected an abnormal peak in your spendings!

		Which account? %s
		Which product? %s
		How much? %f
		How important? %.1f
		When? %s</p>
See more on our website http://re.trackit.io/

this email is originally intended for %s.
`, awsAccount.Pretty, product, an.Cost, an.Cost/an.UpperBand*100, date.String(), ea.Recipient)
	recipient := "team@trackit.io" // replace by ea.Recipient.
	if err := mail.SendMail(recipient, subject, body, ctx); err != nil {
		logger.Error("Error when sending mail", err)
		return err
	} else if err := ea.Insert(db.Db); err != nil {
		logger.Error("Error when updating db with the sent email", err)
		return err
	} else {
		logger.Info("Email sent", map[string]interface{}{
			"aws_account": map[string]interface{}{
				"id":       awsAccount.Id,
				"role_arn": awsAccount.RoleArn,
				"pretty":   awsAccount.Pretty,
				"user_id":  awsAccount.UserId,
			},
			"anomaly": map[string]interface{}{
				"date": an.Date,
				"cost": an.Cost,
			},
			"mail": map[string]interface{}{
				"recipient": recipient,
				"subject":   subject,
				"body":      body,
			},
		})
	}
	return nil
}

// processAnomaliesByProduct checks if the anomaly has already been sent.
// If not, it sends it.
func processAnomaliesByProduct(user users.User, awsAccount aws.AwsAccount, product string, ans []anomalies.CostAnomaly, tx *sql.Tx, ctx context.Context) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	for _, an := range ans {
		if !an.Abnormal {
		} else if date, err := time.Parse("2006-01-02T15:04:05.000Z", an.Date); err != nil {
			logger.Error("Error when parsing date", err)
			return err
		} else if alreadyEmailed, err := models.IsAnomalyAlreadyEmailed(tx, awsAccount.Id, product, date); err != nil {
			logger.Error("Error when checking db for sent email", err)
			return err
		} else if alreadyEmailed {
		} else if err := sendAnomalyEmail(user, awsAccount, product, an, date, ctx); err != nil {
			return nil
		}
	}
	return nil
}

// requestAnomalies requests anomalies from GetAnomaliesData and calls processAnomaliesByProduct.
func requestAnomalies(user users.User, awsAccount aws.AwsAccount, tx *sql.Tx, ctx context.Context) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	identity, err := awsAccount.GetAwsAccountIdentity()
	if err != nil {
		logger.Error("Error when getting identity from AwsAccount", map[string]interface{}{
			"awsAccount": awsAccount,
			"error":      err,
		})
		return err
	}
	params := anomalies.AnomalyEsQueryParams{
		DateBegin:   time.Now().Add(-1 * 24 * 7 * time.Hour),
		DateEnd:     time.Now(),
		AccountList: []string{identity},
	}
	logger.Info("Getting anomalies with following es query params", params)
	anomaliesData, _, err := anomalies.GetAnomaliesData(ctx, params, user)
	if err != nil {
		logger.Error("Error when getting anomalies", err)
		return err
	}
	for product, ans := range anomaliesData {
		if len(product) == 0 {
			continue
		}
		if err := processAnomaliesByProduct(user, awsAccount, product, ans, tx, ctx); err != nil {
			return err
		}
	}
	return nil
}

// prepareAnomaliesRequesting prepares the anomaly checking and calls requestAnomalies.
func prepareAnomaliesRequesting(aas []*models.AwsAccount, tx *sql.Tx, ctx context.Context) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	for _, dbAwsAccount := range aas {
		logger.Info("Processing AWS Account", map[string]interface{}{
			"id":      dbAwsAccount.ID,
			"pretty":  dbAwsAccount.Pretty,
			"user_id": dbAwsAccount.UserID,
		})
		dbUser, err := dbAwsAccount.User(tx)
		if err != nil {
			logger.Error("Error when getting User from AwsAccount", map[string]interface{}{
				"awsAccount": dbAwsAccount,
				"error":      err,
			})
			return err
		}
		awsAccount := aws.AwsAccountFromDbAwsAccount(*dbAwsAccount)
		user := users.UserFromDbUser(*dbUser)
		requestAnomalies(user, awsAccount, tx, ctx)
	}
	return nil
}

// EmailAnomalies emails anomalies detected in the accounts over the past week.
// It doesn't send twice the same email.
// This function should be called at least every week.
func EmailAnomalies(ctx context.Context) (err error) {
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
	var aas []*models.AwsAccount
	if tx, err = db.Db.BeginTx(ctx, nil); err != nil {
	} else if aas, err = models.AwsAccounts(tx); err != nil {
	} else if err = prepareAnomaliesRequesting(aas, tx, ctx); err != nil {
	}
	return
}
