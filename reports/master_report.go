package reports

import (
	"context"
	"github.com/tealeg/xlsx"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit-server/aws"
)

func generateMasterSpreadsheet(ctx context.Context, aa aws.AwsAccount, date string, spreadsheets []*spreadsheet) (*spreadsheet, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	masterSpreadsheet := &spreadsheet{
		aa,
		date,
		xlsx.NewFile(),
	}

//	sheets := make(map[string]sheet, 0)
	for _, spreadsheet := range spreadsheets {
		logger.Debug("-----------TEST--------", nil)
		for _, module := range modules {
			if _, ok := spreadsheet.file.Sheet[module.Name]; !ok {
				logger.Debug("***********", map[string]interface{}{
					"account": spreadsheet.account,
					"sheet": module.Name,
				})
			}
		}
	}

	return masterSpreadsheet, nil
}