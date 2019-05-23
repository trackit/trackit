package reports

import (
	"context"
	"database/sql"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/trackit/trackit-server/aws"
)

type module struct {
	Name          string
	SheetName     string
	ErrorName     string
	GenerateSheet func(context.Context, []aws.AwsAccount, time.Time, *sql.Tx, *excelize.File) error
}

var modules = []module{
	costVariationLastMonth,
	costVariationLast6Months,
	ec2UsageReportModule,
	rdsUsageReportModule,
	esUsageReportModule,
	lambdaUsageReportModule,
	elastiCacheUsageReportModule,
}
