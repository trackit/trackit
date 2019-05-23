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
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
)

/* TODO: Add png into container or store into S3 bucket */
const img = "./reports/introduction.png"

const templateSheetName = "Introduction"

var templateModule = module{
	Name:          "Template",
	SheetName:     templateSheetName,
	ErrorName:     "templateError",
	GenerateSheet: generateTemplateSheet,
}

func generateTemplateSheet(ctx context.Context, _ []aws.AwsAccount, _ time.Time, _ *sql.Tx, file *excelize.File) (err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	file.NewSheet(templateSheetName)
	err = file.AddPicture(templateSheetName, "A1", img, `{"x_scale": 0.95, "y_scale": 1}`)

	if err != nil {
		logger.Error("An error occured while generating template for report", map[string]interface{}{
			"error": err,
		})
	}

	return
}
