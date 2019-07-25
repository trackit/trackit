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
	"database/sql"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/config"
)

const templateSheetName = "Introduction"

var templateModule = module{
	Name:          "Template",
	SheetName:     templateSheetName,
	ErrorName:     "templateError",
	GenerateSheet: generateTemplateSheet,
}

func generateTemplateSheet(ctx context.Context, _ []aws.AwsAccount, _ time.Time, _ *sql.Tx, file *excelize.File) (err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	if len(config.ReportsCover) == 0 {
		return
	}
	imageUrlSplit := strings.Split(config.ReportsCover, "/")
	imageFile := strings.Split(imageUrlSplit[len(imageUrlSplit)-1], ".")
	image, err := downloadFile(config.ReportsCover)
	if err != nil {
		logger.Error("An error occured while downloading cover for report", map[string]interface{}{
			"error": err,
			"cover": config.ReportsCover,
		})
	}
	file.NewSheet(templateSheetName)
	err = file.AddPictureFromBytes(templateSheetName, "A1", `{"x_scale": 0.95, "y_scale": 1}`, imageFile[0], "."+imageFile[len(imageFile)-1], image)
	if err != nil {
		logger.Error("An error occured while generating template for report", map[string]interface{}{
			"error": err,
		})
	}
	return
}
