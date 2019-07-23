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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/tealeg/xlsx"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/awsSession"
	"github.com/trackit/trackit/config"
)

type spreadsheet struct {
	account taws.AwsAccount
	date    string
	file    *xlsx.File
}

type sheet struct {
	name string
	data [][]cell
}

type cell struct {
	value interface{}
	width int
	style *xlsx.Style
}

func (c cell) setValueToCell(newCell *xlsx.Cell) {
	switch value := c.value.(type) {
	case int:
		newCell.SetInt(value)
		break
	case int64:
		newCell.SetInt64(value)
		break
	case float64:
		newCell.SetFloat(value)
		break
	case bool:
		newCell.SetBool(value)
		break
	case string:
		newCell.SetString(value)
		break
	case time.Time:
		newCell.SetDate(value)
		break
	default:
		if defaultValue, ok := value.(string); ok {
			newCell.SetString(defaultValue)
		} else {
			newCell.SetString("Invalid Data")
		}
	}
}

func newCell(value interface{}, dimensions ...int) cell {
	width := 1
	if len(dimensions) > 0 {
		width = dimensions[0]
	}
	item := cell{
		value: value,
		width: width,
		style: xlsx.NewStyle(),
	}
	item.addStyle(defaultStyle{})
	return item
}

func convertToSheet(raw sheet) (sheet xlsx.Sheet) {
	var horizontalPadding int
	sheet = xlsx.Sheet{Name: raw.name}
	for _, rawRow := range raw.data {
		horizontalPadding = 0
		row := sheet.AddRow()
		for _, rawCell := range rawRow {
			for horizontalPadding > 0 {
				row.AddCell()
				horizontalPadding--
			}
			newCell := row.AddCell()
			rawCell.setValueToCell(newCell)
			if rawCell.width > 1 {
				rawCell.width--
				newCell.HMerge = rawCell.width
				horizontalPadding = rawCell.width
			}
			newCell.SetStyle(rawCell.style)
		}
	}
	return
}

func generateSpreadsheet(ctx context.Context, aa taws.AwsAccount, date string, sheets []sheet) (*spreadsheet, map[string]error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Generating spreadsheet", aa)

	file := xlsx.NewFile()
	errors := make(map[string]error)
	for _, rawSheet := range sheets {
		logger.Info("Adding sheet", rawSheet.name)
		sheet := convertToSheet(rawSheet)
		_, err := file.AppendSheet(sheet, rawSheet.name)
		if err != nil {
			logger.Error("Error while adding sheet", map[string]interface{}{
				"sheet": rawSheet.name,
				"error": err.Error(),
			})
			errors[rawSheet.name] = err
		}
	}
	return &spreadsheet{account: aa, date: date, file: file}, errors
}

func getFilenameLocally(account taws.AwsAccount, date string, masterReport bool) string {
	return fmt.Sprintf("/reports/%s", getFilename(account, date, masterReport))
}

func getFilename(account taws.AwsAccount, date string, masterReport bool) string {
	masterReportName := ""
	if masterReport {
		masterReportName = "MasterReport_"
	}
	return fmt.Sprintf("TRACKIT_%s%s_%s.xlsx", masterReportName, account.Pretty, date)
}

func saveSpreadsheetLocally(ctx context.Context, file *spreadsheet, masterReport bool) (err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	filename := getFilenameLocally(file.account, file.date, masterReport)

	err = file.file.Save(filename)
	if err != nil {
		logger.Error("Error while saving file", err)
	}
	return
}

func saveSpreadsheet(ctx context.Context, file *spreadsheet, masterReport bool) (err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	filename := getFilename(file.account, file.date, masterReport)
	reportPath := path.Join(strconv.Itoa(file.account.Id), "generated-report", filename)

	logger.Info("Uploading spreadsheet", reportPath)

	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		err := file.file.Write(writer)
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
