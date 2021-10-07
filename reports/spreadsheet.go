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

package reports

import (
	"context"
	"fmt"
	"io"
	"path"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/awsSession"
	"github.com/trackit/trackit/config"
)

type spreadsheet struct {
	account taws.AwsAccount
	date    string
	File    *excelize.File
}

type spreadsheetType uint8

const (
	MasterReport = spreadsheetType(iota)
	RegularReport
	TagsReport
)

func createSpreadsheet(aa taws.AwsAccount, date string) spreadsheet {
	return spreadsheet{
		account: aa,
		date:    date,
		File:    excelize.NewFile(),
	}
}

/*
func getFilenameLocally(account taws.AwsAccount, date string, reportType spreadsheetType) string {
	return fmt.Sprintf("/reports/%s", getFilename(account, date, reportType))
}
*/

func getFilename(account taws.AwsAccount, date string, reportType spreadsheetType) string {
	reportName := ""
	if reportType == MasterReport {
		reportName = "MasterReport_"
	} else if reportType == TagsReport {
		reportName = "TagsReport_"
	}
	return fmt.Sprintf("TRACKIT_%s%s_%s.xlsx", reportName, account.Pretty, date)
}

/*
func saveSpreadsheetLocally(ctx context.Context, file spreadsheet, reportType spreadsheetType) (err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	filename := getFilenameLocally(file.account, file.date, reportType)

	err = file.File.SaveAs(filename)
	if err != nil {
		logger.Error("Error while saving file", err)
	}
	return
}
*/

func saveSpreadsheet(ctx context.Context, file spreadsheet, reportType spreadsheetType) (err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	filename := getFilename(file.account, file.date, reportType)
	reportPath := path.Join(strconv.Itoa(file.account.Id), "generated-report", filename)

	logger.Info("Uploading spreadsheet", reportPath)

	reader, writer := io.Pipe()

	go func() {
		defer func() {
			if err := writer.Close(); err != nil {
				logger.Error("Error while closing write end of pipe in spreadsheet upload", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
		err := file.File.Write(writer)
		if err != nil {
			logger.Error("Error while saving report", map[string]interface{}{
				"report": reportPath,
				"error":  err.Error(),
			})
		}
	}()

	uploader := s3manager.NewUploader(awsSession.Session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(config.ReportsBucket),
		Key:    aws.String(reportPath),
	})
	if err != nil {
		logger.Error("Failed to upload report", map[string]interface{}{
			"report": reportPath,
			"error":  err.Error(),
		})
	} else {
		logger.Info("Spreadsheet successfully uploaded", result.Location)
	}
	return
}
