package s3

import (
	"context"
	"testing"

	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit2/aws"
)

func init() {
	jsonlog.DefaultLogger = jsonlog.DefaultLogger.WithLogLevel(jsonlog.LogLevelDebug)
}

func TestEpitechio(t *testing.T) {
	err := ReadBills(
		context.Background(),
		taws.AwsAccount{
			RoleArn:  "arn:aws:iam::895365654851:role/trackit",
			External: "RLuxJFYhaZYjWHNYY_pfeAgF@lzymhUKNxiwq_IQ",
		},
		BillRepository{
			Bucket: "epitechio-reports",
			Prefix: "constandusage/",
		},
	)
	if err != nil {
		println(err.Error())
	}
}
